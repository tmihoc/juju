// Copyright 2023 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package bootstrap

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/schema"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	coreapplication "github.com/juju/juju/core/application"
	"github.com/juju/juju/core/arch"
	"github.com/juju/juju/core/base"
	corecharm "github.com/juju/juju/core/charm"
	coreconfig "github.com/juju/juju/core/config"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/objectstore"
	objectstoretesting "github.com/juju/juju/core/objectstore/testing"
	"github.com/juju/juju/core/unit"
	domainapplication "github.com/juju/juju/domain/application"
	applicationcharm "github.com/juju/juju/domain/application/charm"
	applicationservice "github.com/juju/juju/domain/application/service"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/internal/charm"
	"github.com/juju/juju/internal/charm/charmdownloader"
	"github.com/juju/juju/internal/charm/services"
	"github.com/juju/juju/state"
)

type deployerSuite struct {
	baseSuite
}

var _ = gc.Suite(&deployerSuite{})

func (s *deployerSuite) TestValidate(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)
	err := cfg.Validate()
	c.Assert(err, gc.IsNil)

	cfg = s.newConfig(c)
	cfg.DataDir = ""
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.StateBackend = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.CharmUploader = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.ObjectStore = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.ControllerConfig = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.NewCharmRepo = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.NewCharmDownloader = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.CharmhubHTTPClient = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)

	cfg = s.newConfig(c)
	cfg.Logger = nil
	err = cfg.Validate()
	c.Assert(err, jc.ErrorIs, errors.NotValid)
}

func (s *deployerSuite) TestControllerCharmArchWithDefaultArch(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)
	deployer := makeBaseDeployer(cfg)

	arch := deployer.ControllerCharmArch()
	c.Assert(arch, gc.Equals, "amd64")
}

func (s *deployerSuite) TestControllerCharmArch(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.newConfig(c)
	cfg.Constraints = constraints.Value{
		Arch: ptr("arm64"),
	}
	deployer := makeBaseDeployer(cfg)

	arch := deployer.ControllerCharmArch()
	c.Assert(arch, gc.Equals, "arm64")
}

func (s *deployerSuite) TestDeployLocalCharmThatDoesNotExist(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Ensure that if the local charm doesn't exist that we get a not found
	// error. The not found error signals to the caller that they should
	// attempt the charmhub path.

	cfg := s.newConfig(c)
	deployer := makeBaseDeployer(cfg)

	_, err := deployer.DeployLocalCharm(context.Background(), arch.DefaultArchitecture, base.MakeDefaultBase("ubuntu", "22.04"))
	c.Assert(err, jc.ErrorIs, errors.NotFound)
}

func (s *deployerSuite) TestDeployLocalCharm(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// When deploying the local charm, ensure that we get the returned charm URL
	// and the origin with the correct architecture and `/stable` suffix
	// on the origin channel.

	cfg := s.newConfig(c)
	s.ensureControllerCharm(c, cfg.DataDir)

	s.expectLocalCharmUpload(c)

	deployer := s.newBaseDeployer(c, cfg)

	info, err := deployer.DeployLocalCharm(context.Background(), "arm64", base.MakeDefaultBase("ubuntu", "22.04"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(info.URL.String(), gc.Equals, "local:arm64/juju-controller-0")
	c.Assert(info.Origin, gc.DeepEquals, &corecharm.Origin{
		Source: corecharm.Local,
		Type:   "charm",
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04/stable",
		},
	})
	c.Assert(info.Charm, gc.NotNil)
}

func (s *deployerSuite) TestDeployCharmhubCharm(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Ensure we can deploy the charmhub charm, which by default is
	// juju-controller.

	cfg := s.newConfig(c)

	s.expectDownloadAndResolve(c, "juju-controller")

	deployer := s.newBaseDeployer(c, cfg)

	info, err := deployer.DeployCharmhubCharm(context.Background(), "arm64", base.MakeDefaultBase("ubuntu", "22.04"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(info.URL.String(), gc.Equals, "ch:arm64/juju-controller-1")
	c.Assert(info.Origin, gc.DeepEquals, &corecharm.Origin{
		Source:   corecharm.CharmHub,
		Type:     "charm",
		Channel:  &charm.Channel{},
		Hash:     "sha-256",
		Revision: ptr(1),
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	})
	c.Assert(info.Charm, gc.NotNil)
}

func (s *deployerSuite) TestDeployCharmhubCharmWithCustomName(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Ensure we can deploy a charmhub charm with a custom name.

	cfg := s.newConfig(c)
	cfg.ControllerCharmName = "inferi"

	s.expectDownloadAndResolve(c, "inferi")

	deployer := s.newBaseDeployer(c, cfg)

	info, err := deployer.DeployCharmhubCharm(context.Background(), "arm64", base.MakeDefaultBase("ubuntu", "22.04"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(info.URL.String(), gc.Equals, "ch:arm64/inferi-1")
	c.Assert(info.Origin, gc.DeepEquals, &corecharm.Origin{
		Source:   corecharm.CharmHub,
		Type:     "charm",
		Channel:  &charm.Channel{},
		Hash:     "sha-256",
		Revision: ptr(1),
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	})
	c.Assert(info.Charm, gc.NotNil)
}

func (s *deployerSuite) TestAddControllerApplication(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Ensure that we can add the controller application to the model. This will
	// query the backend to ensure that the charm we just uploaded exists before
	// we can add the application.

	cfg := s.newConfig(c)

	curl := "ch:juju-controller-0"

	s.stateBackend.EXPECT().AddApplication(gomock.Any(), s.objectStore).DoAndReturn(func(args state.AddApplicationArgs, store objectstore.ObjectStore) (Application, error) {
		appCfg, err := coreconfig.NewConfig(nil, configSchema, schema.Defaults{
			coreapplication.TrustConfigOptionName: true,
		})
		c.Assert(err, jc.ErrorIsNil)

		// It's interesting that although we don't pass a channel, a stable one
		// is set when persisting the charm origin. I wonder if it would be
		// better to not persist anything at all. In that way we can be sure
		// that we didn't accidentally persist something that we shouldn't have.
		c.Check(args, gc.DeepEquals, state.AddApplicationArgs{
			Name:     bootstrap.ControllerApplicationName,
			Charm:    s.charm,
			CharmURL: "ch:juju-controller-0",
			CharmOrigin: &state.CharmOrigin{
				Source:   "charm-hub",
				Type:     "charm",
				Revision: ptr(1),
				Channel: &state.Channel{
					Risk: "stable",
				},
				Platform: &state.Platform{
					Architecture: "arm64",
					OS:           "ubuntu",
					Channel:      "22.04",
				},
			},
			CharmConfig: map[string]any{
				"is-juju":               true,
				"controller-url":        "wss://obscura.com:1234/api",
				"identity-provider-url": "https://inferi.com",
			},
			Constraints:       constraints.Value{},
			ApplicationConfig: appCfg,
			NumUnits:          1,
		})

		return s.application, nil
	})
	unitName, err := unit.NewNameFromParts(bootstrap.ControllerApplicationName, 0)
	c.Assert(err, jc.ErrorIsNil)
	s.application.EXPECT().Name().Return(bootstrap.ControllerApplicationName)
	s.stateBackend.EXPECT().Unit(unitName.String()).Return(s.unit, nil)
	s.applicationService.EXPECT().CreateApplication(
		gomock.Any(),
		bootstrap.ControllerApplicationName,
		s.charm,
		corecharm.Origin{
			Source:   "charm-hub",
			Type:     "charm",
			Channel:  &charm.Channel{},
			Revision: ptr(1),
			Hash:     "sha-256",
			Platform: corecharm.Platform{
				Architecture: "arm64",
				OS:           "ubuntu",
				Channel:      "22.04",
			},
		},
		applicationservice.AddApplicationArgs{
			ReferenceName: bootstrap.ControllerApplicationName,
			DownloadInfo: &applicationcharm.DownloadInfo{
				CharmhubIdentifier: "abcd",
				Provenance:         applicationcharm.ProvenanceBootstrap,
				DownloadURL:        "https://inferi.com",
				DownloadSize:       42,
			},
			CharmStoragePath: "path",
		},
		applicationservice.AddUnitArg{UnitName: unitName},
	)

	s.charmUploader.EXPECT().PrepareCharmUpload(gomock.Any()).Return(nil, nil)
	s.charmUploader.EXPECT().UpdateUploadedCharm(gomock.Any()).Return(nil, nil)

	deployer := s.newBaseDeployer(c, cfg)

	origin := corecharm.Origin{
		Source:   corecharm.CharmHub,
		Type:     "charm",
		Channel:  &charm.Channel{},
		Revision: ptr(1),
		Hash:     "sha-256",
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	}
	address := "10.0.0.1"
	unit, err := deployer.AddControllerApplication(context.Background(), DeployCharmInfo{
		URL:    charm.MustParseURL(curl),
		Charm:  s.charm,
		Origin: &origin,
		DownloadInfo: &corecharm.DownloadInfo{
			CharmhubIdentifier: "abcd",
			DownloadURL:        "https://inferi.com",
			DownloadSize:       42,
		},
		ArchivePath:     "path",
		ObjectStoreUUID: "1234",
	}, address)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unit, gc.NotNil)
}

func (s *deployerSuite) ensureControllerCharm(c *gc.C, dataDir string) {
	// This will place the most basic charm (no hooks, no config, no actions)
	// into the data dir so that we can test the local charm path.

	metadata := `
name: juju-controller
summary: Juju controller
description: Juju controller
`
	dir := c.MkDir()
	err := os.WriteFile(filepath.Join(dir, "metadata.yaml"), []byte(metadata), 0644)
	c.Assert(err, jc.ErrorIsNil)

	charmDir, err := charm.ReadCharmDir(dir)
	c.Assert(err, jc.ErrorIsNil)

	var buf bytes.Buffer
	charmDir.ArchiveTo(&buf)

	path := filepath.Join(dataDir, "charms")
	err = os.MkdirAll(path, 0755)
	c.Assert(err, jc.ErrorIsNil)

	path = filepath.Join(path, bootstrap.ControllerCharmArchive)
	err = os.WriteFile(path, buf.Bytes(), 0644)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *deployerSuite) newBaseDeployer(c *gc.C, cfg BaseDeployerConfig) baseDeployer {
	deployer := makeBaseDeployer(cfg)

	deployer.stateBackend = s.stateBackend
	deployer.charmUploader = s.charmUploader
	deployer.objectStore = s.objectStore

	return deployer
}

func (s *deployerSuite) expectLocalCharmUpload(c *gc.C) {
	s.charmUploader.EXPECT().PrepareLocalCharmUpload("local:juju-controller-0").Return(&charm.URL{
		Schema:       "local",
		Name:         "juju-controller",
		Revision:     0,
		Architecture: "arm64",
	}, nil)
	// Ensure that the charm uploaded to the object store is the one we expect.
	s.objectStore.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, path string, reader io.Reader, size int64) (objectstore.UUID, error) {
		c.Check(strings.HasPrefix(path, "charms/local:arm64/juju-controller-0"), jc.IsTrue)
		return "", nil
	})
	s.charmUploader.EXPECT().UpdateUploadedCharm(gomock.Any()).DoAndReturn(func(info state.CharmInfo) (services.UploadedCharm, error) {
		c.Check(info.ID, gc.Equals, "local:arm64/juju-controller-0")
		c.Check(strings.HasPrefix(info.StoragePath, "charms/local:arm64/juju-controller-0"), jc.IsTrue)

		return nil, nil
	})
}

func (s *deployerSuite) expectDownloadAndResolve(c *gc.C, name string) {
	curl := &charm.URL{
		Schema:       string(charm.CharmHub),
		Name:         name,
		Revision:     0,
		Architecture: "arm64",
	}
	origin := corecharm.Origin{
		Source:  corecharm.CharmHub,
		Type:    "charm",
		Channel: &charm.Channel{},
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	}
	resolvedOrigin := corecharm.Origin{
		Source:   corecharm.CharmHub,
		Type:     "charm",
		Channel:  &charm.Channel{},
		Hash:     "sha-256",
		Revision: ptr(1),
		Platform: corecharm.Platform{
			Architecture: "arm64",
			OS:           "ubuntu",
			Channel:      "22.04",
		},
	}

	s.charmRepo.EXPECT().ResolveWithPreferredChannel(gomock.Any(), name, origin).Return(corecharm.ResolvedData{
		URL:    curl,
		Origin: resolvedOrigin,
		EssentialMetadata: corecharm.EssentialMetadata{
			ResolvedOrigin: resolvedOrigin,
			DownloadInfo: corecharm.DownloadInfo{
				DownloadURL: "https://inferi.com",
			},
		},
	}, nil)

	url, err := url.Parse("https://inferi.com")
	c.Assert(err, jc.ErrorIsNil)

	s.charmDownloader.EXPECT().Download(gomock.Any(), url, "sha-256").Return(&charmdownloader.DownloadResult{
		SHA256: "sha-256",
		SHA384: "sha-384",
		Path:   "path",
		Size:   42,
	}, nil)

	objectStoreUUID := objectstoretesting.GenObjectStoreUUID(c)

	s.applicationService.EXPECT().ResolveControllerCharmDownload(gomock.Any(), domainapplication.ResolveControllerCharmDownload{
		SHA256: "sha-256",
		SHA384: "sha-384",
		Path:   "path",
		Size:   42,
	}).Return(domainapplication.ResolvedControllerCharmDownload{
		Charm:           s.charm,
		ArchivePath:     "path",
		ObjectStoreUUID: objectStoreUUID,
	}, nil)
}
