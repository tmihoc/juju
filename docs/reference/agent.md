---
myst:
  html_meta:
    description: "Understand Juju agents: the controller, machine and unit agents, and the model workers, that manage deployments through jujuagentd and containeragent processes."
---

(agent)=
# Agent

In Juju, an **agent** is a {ref}`jujuagentd` / {ref}`containeragent` process that works to realise the state declared by a Juju end-user with a Juju client (e.g., {ref}`the juju CLI <juju-cli>`) for a Juju entity (e.g., {ref}`controller <controller>`, {ref}`model <model>`, {ref}`machine <machine>`, {ref}`unit <unit>`) via [workers](https://pkg.go.dev/github.com/juju/worker/v5).

On machines, an agent is managed by `systemd`.

## Types of agents

- **Controller agent:** The agent that runs a {ref}`controller <controller>`. See {ref}`controller-agent`.
- **Machine agent:** The agent that runs on a {ref}`machine <machine>`. See {ref}`machine-agent`.
- **Model workers:** The workers that act for a {ref}`model <model>` run inside the controller agent, with a model operator pod for a Kubernetes model. See {ref}`model-workers`.
- **Unit agent:** The agent of a {ref}`unit <unit>`, described next.

(unit-agent)=
## Unit agent

The unit agent runs a {ref}`unit <unit>`: its `uniter` worker installs the charm and runs its hooks. On machine clouds the machine agent's deployer starts it as an engine of workers inside the machine agent's own process, and on Kubernetes it is a `containeragent` process in the unit's pod.

## How agents learn of changes

When a Juju user uses the client (e.g., types a command in the CLI), this goes to the controller agent's `apiserver`, which passes it on to the database. The database runs a background process that checks if anything has changed and, if so, emits an event (think "I've seen something that's changed. Do you care about it?"). The event cascades through Juju. The unit agent becomes aware of it by always polling the controller agent as part of a reconciliation loop trying to reconcile the unit agent's local state to the remote state on the controller (i.e., the state in the controller's database).

