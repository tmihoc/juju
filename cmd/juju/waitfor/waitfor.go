// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package waitfor

import (
	"github.com/juju/cmd/v3"

	_ "github.com/juju/juju/provider/all"
)

// Logger is the interface used by the wait-for command to log messages.
type Logger interface {
	Infof(string, ...any)
	Verbosef(string, ...any)
}

var waitForDoc = `
The wait-for set of commands (model, application, machine and unit) represents 
a way to wait for a goal state to be reached. The goal state can be defined
programmatically using the query DSL (domain specific language).

The query DSL is a simple language that can be comprised of expressions to
produce a boolean result. The result of the query is used to determine if the
goal state has been reached. The query DSL is evaluated against the scope of
the command.

Built-in functions are provided to help define the goal state. The built-in
functions are defined in the query package. Examples of built-in functions
include len, print, forEach (lambda), startsWith and endsWith.

Examples:
    juju wait-for unit mysql/0
    juju wait-for application mysql --query='name=="mysql" && (status=="active" || status=="idle")'
    juju wait-for model default --query='forEach(units, unit => startsWith(unit.name, "ubuntu"))'

See also:
    wait-for model
    wait-for application
    wait-for machine
    wait-for unit
`

// NewWaitForCommand creates the wait-for supercommand and registers the
// subcommands that it supports.
func NewWaitForCommand() cmd.Command {
	waitFor := cmd.NewSuperCommand(cmd.SuperCommandParams{
		Name:        "wait-for",
		UsagePrefix: "juju",
		Doc:         waitForDoc,
		Purpose:     "Wait for an entity to reach a specified state."})

	waitFor.Register(newApplicationCommand())
	waitFor.Register(newMachineCommand())
	waitFor.Register(newModelCommand())
	waitFor.Register(newUnitCommand())
	return waitFor
}
