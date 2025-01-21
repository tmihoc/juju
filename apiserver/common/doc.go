// Copyright 2025 Canonical Ltd. Licensed under the AGPLv3, see LICENCE file for
// details.

// For facades that need to share functionality, this functionality is located
// in package common. For example, multiple facades (e.g., deployer,
// storageprovisioner, instancemutator, etc.) need to track entity life, and
// this shared functionality is abstracted away in apiserver/common/life.go.
//
// Functionality specific to a single facade should be kept in that facade's
// directory.
//
// We choose to use a 'common' directory to avoid problems with circular
// dependencies. Otherwise it is very easy for facade A to depend on some
// functionality from facade B, but we end up with B then needing something that
// A provided. Also, pulling things into 'common' makes it clearer that this
// functionality *is* shared, and thus care needs to be taken that any changes
// made will impact more than one facade.

package common
