// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package agentbootstrap_test

import (
	"context"

	"github.com/juju/errors"
	mgotesting "github.com/juju/mgo/v3/testing"
	"github.com/juju/names/v6"
	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/agent/agentbootstrap"
	"github.com/juju/juju/cloud"
	"github.com/juju/juju/controller"
	corebase "github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/model"
	corenetwork "github.com/juju/juju/core/network"
	jujuversion "github.com/juju/juju/core/version"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/envcontext"
	"github.com/juju/juju/internal/charmhub"
	"github.com/juju/juju/internal/cloudconfig/instancecfg"
	"github.com/juju/juju/internal/database"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/mongo"
	"github.com/juju/juju/internal/mongo/mongotest"
	"github.com/juju/juju/internal/network"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/storage/provider"
	"github.com/juju/juju/internal/testing"
	jujujujutesting "github.com/juju/juju/juju/testing"
	"github.com/juju/juju/state"
)

type bootstrapSuite struct {
	testing.BaseSuite
	mgoInst mgotesting.MgoInstance
}

var _ = gc.Suite(&bootstrapSuite{})

func (s *bootstrapSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	// Don't use MgoSuite, because we need to ensure
	// we have a fresh mongo for each test case.
	s.mgoInst.EnableAuth = true
	s.mgoInst.EnableReplicaSet = true
	err := s.mgoInst.Start(testing.Certs)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *bootstrapSuite) TearDownTest(c *gc.C) {
	s.mgoInst.Destroy()
	s.BaseSuite.TearDownTest(c)
}

func (s *bootstrapSuite) TestInitializeState(c *gc.C) {
	dataDir := c.MkDir()

	s.PatchValue(&network.AddressesForInterfaceName, func(name string) ([]string, error) {
		if name == network.DefaultLXDBridge {
			return []string{
				"10.0.4.1",
				"10.0.4.4",
			}, nil
		}
		c.Fatalf("unknown bridge in testing: %v", name)
		return nil, nil
	})

	configParams := agent.AgentConfigParams{
		Paths:             agent.Paths{DataDir: dataDir},
		Tag:               names.NewMachineTag("0"),
		UpgradedToVersion: jujuversion.Current,
		APIAddresses:      []string{"localhost:17070"},
		CACert:            testing.CACert,
		Password:          testing.DefaultMongoPassword,
		Controller:        testing.ControllerTag,
		Model:             testing.ModelTag,
	}
	servingInfo := controller.StateServingInfo{
		Cert:           testing.ServerCert,
		PrivateKey:     testing.ServerKey,
		CAPrivateKey:   testing.CAKey,
		APIPort:        1234,
		StatePort:      s.mgoInst.Port(),
		SystemIdentity: "def456",
	}

	cfg, err := agent.NewStateMachineConfig(configParams, servingInfo)
	c.Assert(err, jc.ErrorIsNil)

	_, available := cfg.StateServingInfo()
	c.Assert(available, jc.IsTrue)
	expectBootstrapConstraints := constraints.MustParse("mem=1024M")
	expectModelConstraints := constraints.MustParse("mem=512M")
	initialAddrs := corenetwork.NewMachineAddresses([]string{
		"zeroonetwothree",
		"0.1.2.3",
		"10.0.3.3", // not a lxc bridge address
		"10.0.4.1", // lxd bridge address filtered.
		"10.0.4.4", // lxd bridge address filtered.
		"10.0.4.5", // not a lxd bridge address
	}).AsProviderAddresses()

	modelAttrs := testing.FakeConfig().Merge(testing.Attrs{
		"agent-version":  jujuversion.Current.String(),
		"charmhub-url":   charmhub.DefaultServerURL,
		"not-for-hosted": "foo",
	})
	modelCfg, err := config.New(config.NoDefaults, modelAttrs)
	c.Assert(err, jc.ErrorIsNil)
	controllerCfg := testing.FakeControllerConfig()

	controllerInheritedConfig := map[string]interface{}{
		"apt-mirror": "http://mirror",
		"no-proxy":   "value",
	}
	regionConfig := cloud.RegionConfig{
		"some-region": cloud.Attrs{
			"no-proxy": "a-value",
		},
	}
	registry := provider.CommonStorageProviders()
	var envProvider fakeProvider
	stateInitParams := instancecfg.StateInitializationParams{
		BootstrapMachineConstraints: expectBootstrapConstraints,
		BootstrapMachineInstanceId:  "i-bootstrap",
		BootstrapMachineDisplayName: "test-display-name",
		ControllerCloud: cloud.Cloud{
			Name:         "dummy",
			Type:         "dummy",
			AuthTypes:    []cloud.AuthType{cloud.EmptyAuthType},
			Regions:      []cloud.Region{{Name: "dummy-region"}},
			RegionConfig: regionConfig,
		},
		ControllerCloudRegion:         "dummy-region",
		ControllerConfig:              controllerCfg,
		ControllerModelConfig:         modelCfg,
		ControllerModelEnvironVersion: 666,
		ModelConstraints:              expectModelConstraints,
		ControllerInheritedConfig:     controllerInheritedConfig,
		StoragePools: map[string]storage.Attrs{
			"spool": {
				"type": "loop",
				"foo":  "bar",
			},
		},
	}
	adminUser := names.NewLocalUserTag("agent-admin")
	bootstrap, err := agentbootstrap.NewAgentBootstrap(
		agentbootstrap.AgentBootstrapArgs{
			AgentConfig:               cfg,
			BootstrapEnviron:          &fakeEnviron{},
			AdminUser:                 adminUser,
			StateInitializationParams: stateInitParams,
			MongoDialOpts:             mongotest.DialOpts(),
			BootstrapMachineAddresses: initialAddrs,
			BootstrapMachineJobs:      []model.MachineJob{model.JobManageModel},
			SharedSecret:              "abc123",
			StorageProviderRegistry:   registry,
			BootstrapDqlite:           bootstrapDqliteWithDummyCloudType,
			Provider: func(t string) (environs.EnvironProvider, error) {
				c.Assert(t, gc.Equals, "dummy")
				return &envProvider, nil
			},
			Logger: loggertesting.WrapCheckLog(c),
		},
	)
	c.Assert(err, jc.ErrorIsNil)

	ctlr, err := bootstrap.Initialize(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	defer func() { _ = ctlr.Close() }()

	st, err := ctlr.SystemState()
	c.Assert(err, jc.ErrorIsNil)
	err = cfg.Write()
	c.Assert(err, jc.ErrorIsNil)

	// Check that the model has been set up.
	model, err := st.Model()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(model.UUID(), gc.Equals, modelCfg.UUID())
	c.Check(model.EnvironVersion(), gc.Equals, 666)

	// Check that initial admin user has been set up correctly.
	modelTag := model.Tag().(names.ModelTag)
	controllerTag := names.NewControllerTag(controllerCfg.ControllerUUID())
	s.assertCanLogInAsAdmin(c, modelTag, controllerTag, testing.DefaultMongoPassword)

	// Check that controller model configuration has been added, and
	// model constraints set.
	model, err = st.Model()
	c.Assert(err, jc.ErrorIsNil)

	gotModelConstraints, err := st.ModelConstraints()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(gotModelConstraints, gc.DeepEquals, expectModelConstraints)

	// Check that the bootstrap machine looks correct.
	m, err := st.Machine("0")
	c.Assert(err, jc.ErrorIsNil)
	c.Check(m.Id(), gc.Equals, "0")
	c.Check(m.Jobs(), gc.DeepEquals, []state.MachineJob{state.JobManageModel})

	base, err := corebase.ParseBase(m.Base().OS, m.Base().Channel)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(m.Base().String(), gc.Equals, base.String())
	c.Check(m.CheckProvisioned(agent.BootstrapNonce), jc.IsTrue)

	gotBootstrapConstraints, err := m.Constraints()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(gotBootstrapConstraints, gc.DeepEquals, expectBootstrapConstraints)

	// Check that the state serving info is initialised correctly.
	stateServingInfo, err := st.StateServingInfo()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(stateServingInfo, jc.DeepEquals, controller.StateServingInfo{
		APIPort:        1234,
		StatePort:      s.mgoInst.Port(),
		Cert:           testing.ServerCert,
		PrivateKey:     testing.ServerKey,
		CAPrivateKey:   testing.CAKey,
		SharedSecret:   "abc123",
		SystemIdentity: "def456",
	})

	// Check that the machine agent's config has been written
	// and that we can use it to connect to mongo.
	machine0 := names.NewMachineTag("0")
	newCfg, err := agent.ReadConfig(agent.ConfigPath(dataDir, machine0))
	c.Assert(err, jc.ErrorIsNil)
	c.Check(newCfg.Tag(), gc.Equals, machine0)

	info, ok := cfg.MongoInfo()
	c.Assert(ok, jc.IsTrue)
	c.Check(info.Password, gc.Not(gc.Equals), testing.DefaultMongoPassword)

	session, err := mongo.DialWithInfo(*info, mongotest.DialOpts())
	c.Assert(err, jc.ErrorIsNil)
	session.Close()
}

func (s *bootstrapSuite) TestInitializeStateWithStateServingInfoNotAvailable(c *gc.C) {
	configParams := agent.AgentConfigParams{
		Paths:             agent.Paths{DataDir: c.MkDir()},
		Tag:               names.NewMachineTag("0"),
		UpgradedToVersion: jujuversion.Current,
		APIAddresses:      []string{"localhost:17070"},
		CACert:            testing.CACert,
		Password:          "fake",
		Controller:        testing.ControllerTag,
		Model:             testing.ModelTag,
	}
	cfg, err := agent.NewAgentConfig(configParams)
	c.Assert(err, jc.ErrorIsNil)

	_, available := cfg.StateServingInfo()
	c.Assert(available, jc.IsFalse)

	adminUser := names.NewLocalUserTag("agent-admin")

	bootstrap, err := agentbootstrap.NewAgentBootstrap(
		agentbootstrap.AgentBootstrapArgs{
			AgentConfig:               cfg,
			BootstrapEnviron:          &fakeEnviron{},
			AdminUser:                 adminUser,
			StateInitializationParams: instancecfg.StateInitializationParams{},
			MongoDialOpts:             mongotest.DialOpts(),
			SharedSecret:              "abc123",
			StorageProviderRegistry:   provider.CommonStorageProviders(),
			BootstrapDqlite:           bootstrapDqliteWithDummyCloudType,
			Logger:                    loggertesting.WrapCheckLog(c),
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	_, err = bootstrap.Initialize(context.Background())

	// InitializeState will fail attempting to get the api port information
	c.Assert(err, gc.ErrorMatches, "state serving information not available")
}

func (s *bootstrapSuite) TestInitializeStateFailsSecondTime(c *gc.C) {
	dataDir := c.MkDir()

	configParams := agent.AgentConfigParams{
		Paths:             agent.Paths{DataDir: dataDir},
		Tag:               names.NewMachineTag("0"),
		UpgradedToVersion: jujuversion.Current,
		APIAddresses:      []string{"localhost:17070"},
		CACert:            testing.CACert,
		Password:          testing.DefaultMongoPassword,
		Controller:        testing.ControllerTag,
		Model:             testing.ModelTag,
	}
	cfg, err := agent.NewAgentConfig(configParams)
	c.Assert(err, jc.ErrorIsNil)
	cfg.SetStateServingInfo(controller.StateServingInfo{
		APIPort:        5555,
		StatePort:      s.mgoInst.Port(),
		Cert:           testing.CACert,
		PrivateKey:     testing.CAKey,
		SharedSecret:   "baz",
		SystemIdentity: "qux",
	})
	modelAttrs := testing.FakeConfig().Delete("admin-secret").Merge(testing.Attrs{
		"agent-version": jujuversion.Current.String(),
		"charmhub-url":  charmhub.DefaultServerURL,
	})
	modelCfg, err := config.New(config.NoDefaults, modelAttrs)
	c.Assert(err, jc.ErrorIsNil)

	args := instancecfg.StateInitializationParams{
		BootstrapMachineInstanceId:  "i-bootstrap",
		BootstrapMachineDisplayName: "test-display-name",
		ControllerCloud: cloud.Cloud{
			Name:      "dummy",
			Type:      "dummy",
			AuthTypes: []cloud.AuthType{cloud.EmptyAuthType},
			Regions:   []cloud.Region{{Name: "dummy-region"}},
		},
		ControllerConfig:      testing.FakeControllerConfig(),
		ControllerModelConfig: modelCfg,
	}

	adminUser := names.NewLocalUserTag("agent-admin")
	bootstrap, err := agentbootstrap.NewAgentBootstrap(
		agentbootstrap.AgentBootstrapArgs{
			AgentConfig:               cfg,
			BootstrapEnviron:          &fakeEnviron{},
			AdminUser:                 adminUser,
			StateInitializationParams: args,
			MongoDialOpts:             mongotest.DialOpts(),
			BootstrapMachineJobs:      []model.MachineJob{model.JobManageModel},
			SharedSecret:              "abc123",
			StorageProviderRegistry:   provider.CommonStorageProviders(),
			BootstrapDqlite:           bootstrapDqliteWithDummyCloudType,
			Provider: func(t string) (environs.EnvironProvider, error) {
				return &fakeProvider{}, nil
			},
			Logger: loggertesting.WrapCheckLog(c),
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	st, err := bootstrap.Initialize(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	_ = st.Close()

	bootstrap, err = agentbootstrap.NewAgentBootstrap(
		agentbootstrap.AgentBootstrapArgs{
			AgentConfig:               cfg,
			BootstrapEnviron:          &fakeEnviron{},
			AdminUser:                 adminUser,
			StateInitializationParams: args,
			MongoDialOpts:             mongotest.DialOpts(),
			SharedSecret:              "baz",
			StorageProviderRegistry:   provider.CommonStorageProviders(),
			BootstrapDqlite:           database.BootstrapDqlite,
			Logger:                    loggertesting.WrapCheckLog(c),
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	st, err = bootstrap.Initialize(context.Background())
	if err == nil {
		_ = st.Close()
	}
	c.Assert(err, jc.ErrorIs, errors.AlreadyExists)
}

func (s *bootstrapSuite) TestMachineJobFromParams(c *gc.C) {
	var tests = []struct {
		name model.MachineJob
		want state.MachineJob
		err  string
	}{{
		name: model.JobHostUnits,
		want: state.JobHostUnits,
	}, {
		name: model.JobManageModel,
		want: state.JobManageModel,
	}, {
		name: "invalid",
		want: -1,
		err:  `invalid machine job "invalid"`,
	}}
	for _, test := range tests {
		got, err := agentbootstrap.MachineJobFromParams(test.name)
		if err != nil {
			c.Check(err, gc.ErrorMatches, test.err)
		}
		c.Check(got, gc.Equals, test.want)
	}
}

func (s *bootstrapSuite) assertCanLogInAsAdmin(c *gc.C, modelTag names.ModelTag, controllerTag names.ControllerTag, password string) {
	session, err := mongo.DialWithInfo(mongo.MongoInfo{
		Info: mongo.Info{
			Addrs:  []string{s.mgoInst.Addr()},
			CACert: testing.CACert,
		},
		Tag:      nil, // admin user
		Password: password,
	}, mongotest.DialOpts())
	c.Assert(err, jc.ErrorIsNil)
	session.Close()
}

type fakeProvider struct {
	environs.EnvironProvider
	jujutesting.Stub
	environ *fakeEnviron
}

func (p *fakeProvider) ValidateCloud(_ context.Context, spec cloudspec.CloudSpec) error {
	p.MethodCall(p, "ValidateCloud", spec)
	return p.NextErr()
}

func (p *fakeProvider) Validate(_ context.Context, newCfg, oldCfg *config.Config) (*config.Config, error) {
	p.MethodCall(p, "Validate", newCfg, oldCfg)
	return newCfg, p.NextErr()
}

func (p *fakeProvider) Open(_ context.Context, args environs.OpenParams) (environs.Environ, error) {
	p.MethodCall(p, "Open", args)
	p.environ = &fakeEnviron{Stub: &p.Stub, provider: p}
	return p.environ, p.NextErr()
}

func (p *fakeProvider) Version() int {
	p.MethodCall(p, "Version")
	p.PopNoErr()
	return 123
}

type fakeEnviron struct {
	environs.Environ
	*jujutesting.Stub
	provider *fakeProvider

	callCtxUsed envcontext.ProviderCallContext
}

func (e *fakeEnviron) Create(ctx envcontext.ProviderCallContext, args environs.CreateParams) error {
	e.MethodCall(e, "Create", ctx, args)
	e.callCtxUsed = ctx
	return e.NextErr()
}

func (e *fakeEnviron) Provider() environs.EnvironProvider {
	e.MethodCall(e, "Provider")
	e.PopNoErr()
	return e.provider
}

func bootstrapDqliteWithDummyCloudType(
	ctx context.Context,
	mgr database.BootstrapNodeManager,
	modelUUID model.UUID,
	logger logger.Logger,
	opts ...database.BootstrapOpt,
) error {
	// The dummy cloud type needs to be inserted before the other operations.
	opts = append([]database.BootstrapOpt{
		jujujujutesting.InsertDummyCloudType,
	}, opts...)

	return database.BootstrapDqlite(ctx, mgr, modelUUID, logger, opts...)
}
