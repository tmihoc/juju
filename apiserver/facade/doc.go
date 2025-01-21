// Copyright 2025 Canonical Ltd. Licensed under the AGPLv3, see LICENCE file for
// details.

// Package facade implements the registry for facades (file registry.go) and the
// context that is passed into facades when they’re instantiated (file
// interface.go).
//
// Facades add a function to the registry by calling registry.MustRegister,
// passing in a factory function that will take a facade.Context and return a
// facade.Facade.
//
// The registry is used by the API server to track what functionality is
// available, and to execute that functionality on a client request.

package facade
