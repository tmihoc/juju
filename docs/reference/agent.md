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

Agents nest. A controller agent starts an engine of model workers for each model, and a machine agent starts the unit agents of the units on its machine. On a controller machine the machine agent is the controller agent. Each type is described with the entity it runs.

### Controller agent

The `jujuagentd` process of a controller node: it runs the API server, the database and the other workers only a controller needs. See {ref}`controller-agent` on the {ref}`controller <controller>` page.

### Model workers

The workers that act for one {ref}`model <model>`. They run as an engine inside the controller agent, one engine per activated model, and a Kubernetes model also has a model operator pod. See {ref}`model-workers`.

### Machine agent

The `jujuagentd` process on a {ref}`machine <machine>`: it reports the machine's health, runs its containers and its unit agents, and marks the machine dead when it is removed. See {ref}`machine-agent`.

(unit-agent)=
### Unit agent

The workers of a {ref}`unit <unit>`: its `uniter` worker installs the charm and runs its hooks. On machine clouds the machine agent's deployer starts it as an engine inside the machine agent's own process, and on Kubernetes it is a `containeragent` process in the unit's pod.

## How agents learn of changes

When a Juju user uses the client (e.g., types a command in the CLI), this goes to the controller agent's `apiserver`, which passes it on to the database. The database runs a background process that checks if anything has changed and, if so, emits an event (think "I've seen something that's changed. Do you care about it?"). The event cascades through Juju. The unit agent becomes aware of it by always polling the controller agent as part of a reconciliation loop trying to reconcile the unit agent's local state to the remote state on the controller (i.e., the state in the controller's database).

