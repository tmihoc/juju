---
myst:
  html_meta:
    description: "Juju model reference: logical workspaces for applications, resources, and configuration. Learn controller and workload model types."
---

(model)=
# Model
```{ibnote}
See also: {ref}`manage-models`
```

In Juju, a **model**  is an abstraction that holds {ref}`applications <application>` and application supporting components -- {ref}`machines <machine>`, {ref}`storage <storage>`, {ref}`network spaces <space>`, {ref}`relations <relation>`, etc.

A model is created by a {ref}`user <user>`, and owned in perpetuity by that user (or a new user with the same name), though it may also be used by any other user with model access level, within the limits of their level.

A model is created on a {ref}`controller <controller>`.  Both the model and the controller are associated with a {ref}`cloud <cloud>` (and a cloud {ref}`credential <credential>`), though they do not both have to be on the same cloud (this is a scenario where you have a 'multicloud controller' and where you may have 'cross-model relations'). Any entities added to the model will use resources from that cloud.

One can deploy multiple applications to the same model. Thus, models allow the logical grouping of applications and infrastructure that work together to deliver a service or product.  Moreover, one can apply common {ref}`configurations <configuration>` to a whole model. As such, models allow the low-level storage, compute, network and software components to be reasoned about as a single entity as well.

## Model taxonomy

Models are of two types:

1. **The controller model (`controller`).** This is your Juju management model. A Juju deployment will have just one controller model, which is created by default when you create a controller (`juju bootstrap`). It typically contains a single machine, for the controller (the `controller` application). If controller {ref}`high availability <high-availability>` is enabled, then the controller model would contain multiple instances. The `controller` model may also contain certain applications which it makes sense to deploy near the controller -- e.g., starting with Juju `3.0`, the `juju-dashboard` application.

2. **Regular model.** This is your Juju workload model. A Juju deployment may have many different workload models, which you create manually (`juju add-model`). It is the model where you typically deploy your applications.

## Model configuration

A model configuration is a rule or a set of rules that define the behavior of a model -- including the `controller` model.

```{ibnote}
See more: {ref}`list-of-model-configuration-keys`,  {ref}`configure-a-model`
```

## Model data model

```{ggarch}
:file: ../juju.ggarch
:view: Model data model
:caption: The model database stores a model in a `model` record: its uuid, name, qualifier (the owner namespace of the name), type, cloud, cloud type and region, the credential that the model uses, whether it is the controller model, and the uuid of its controller. Three records point at the model by its uuid: `model_life` holds its life (the stage it is in: alive, dying, or dead), `model_constraint` points at its constraint, and `model_agent` holds the password hash of the model's agent and the algorithm that made it.
:alt: An entity-relationship diagram of four tables. The model table is in the middle, with model_constraint on its left, model_agent on its right, and model_life below it. Each of the three has a line to the model table through its model_uuid column.
```

## Model migration

```{ggarch}
:file: ../juju.ggarch
:view: Model migration phase
:caption: A model migration (moving a model from one controller to another, started with `juju migrate`) is driven by the migration master (the worker on the source controller that runs the migration) and always moves through phases in the order shown. In QUIESCE the source controller runs prechecks and waits for the model's agents to stop work. IMPORT transfers the model's data and binaries to the target controller. In VALIDATION the agents check the model against the target, the source checks the machines in the migrated model, and the model is activated on the target. In SUCCESS the agents report that they run against the target and the model's resources move across. LOGTRANSFER copies the model's logs. REAP removes the model from the source and ends in DONE, or in REAPFAILED if the removal fails. A failure in QUIESCE, IMPORT, or VALIDATION (a failed precheck, transfer, validation, machine check, or activation) moves the migration to ABORT, which removes the imported model from the target and ends in ABORTDONE. SUCCESS and later phases have no abort. DONE, REAPFAILED, and ABORTDONE are the terminal phases.
:alt: A chain of phases from QUIESCE through IMPORT, VALIDATION, SUCCESS, LOGTRANSFER, and REAP, which splits to DONE and REAPFAILED. Arrows labelled abort lead from QUIESCE, IMPORT, and VALIDATION to ABORT, which leads to ABORTDONE.
```
