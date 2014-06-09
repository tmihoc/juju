// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"os"

	"github.com/juju/testing"
	"github.com/juju/utils"
	gc "launchpad.net/gocheck"

	"github.com/juju/juju/juju/osenv"
)

// JujuOsEnvSuite isolates the tests from Juju environment variables.
// This is intended to be only used by existing suites, usually embedded in
// BaseSuite and in FakeJujuHomeSuite. Eventually the tests relying on
// JujuOsEnvSuite will be converted to use the IsolationSuite in
// github.com/juju/testing, and this suite will be removed.
// Do not use JujuOsEnvSuite when writing new tests.
type JujuOsEnvSuite struct {
	oldHomeEnv     string
	oldEnvironment map[string]string
}

func (s *JujuOsEnvSuite) SetUpSuite(c *gc.C) {
}

func (s *JujuOsEnvSuite) TearDownSuite(c *gc.C) {
}

func (s *JujuOsEnvSuite) SetUpTest(c *gc.C) {
	s.oldEnvironment = make(map[string]string)
	for _, name := range []string{
		osenv.JujuHomeEnvKey,
		osenv.JujuEnvEnvKey,
		osenv.JujuLoggingConfigEnvKey,
	} {
		s.oldEnvironment[name] = os.Getenv(name)
		os.Setenv(name, "")
	}
	s.oldHomeEnv = utils.Home()
	utils.SetHome("")
}

func (s *JujuOsEnvSuite) TearDownTest(c *gc.C) {
	for name, value := range s.oldEnvironment {
		os.Setenv(name, value)
	}
	utils.SetHome(s.oldHomeEnv)
}

// BaseSuite provides required functionality for all test suites
// when embedded in a gocheck suite type:
// - logger redirect
// - no outgoing network access
// - protection of user's home directory
// - scrubbing of env vars
// TODO (frankban) 2014-06-09: switch to using IsolationSuite.
type BaseSuite struct {
	testing.CleanupSuite
	testing.LoggingSuite
	JujuOsEnvSuite
}

func (s *BaseSuite) SetUpSuite(c *gc.C) {
	s.CleanupSuite.SetUpSuite(c)
	s.LoggingSuite.SetUpSuite(c)
	s.JujuOsEnvSuite.SetUpSuite(c)
	s.PatchValue(&utils.OutgoingAccessAllowed, false)
}

func (s *BaseSuite) TearDownSuite(c *gc.C) {
	s.JujuOsEnvSuite.TearDownSuite(c)
	s.LoggingSuite.TearDownSuite(c)
	s.CleanupSuite.TearDownSuite(c)
}

func (s *BaseSuite) SetUpTest(c *gc.C) {
	s.CleanupSuite.SetUpTest(c)
	s.LoggingSuite.SetUpTest(c)
	s.JujuOsEnvSuite.SetUpTest(c)
}

func (s *BaseSuite) TearDownTest(c *gc.C) {
	s.JujuOsEnvSuite.TearDownTest(c)
	s.LoggingSuite.TearDownTest(c)
	s.CleanupSuite.TearDownTest(c)
}
