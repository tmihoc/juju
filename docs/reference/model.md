---
myst:
  html_meta:
    description: "Juju model reference: logical workspaces for applications, resources, and configuration. The workers that act for a model, the records Juju stores about it, the change streams it publishes, and the operations that create, configure, migrate and remove models."
---

(model)=
# Model

In Juju, a **model** is an abstraction that holds {ref}`applications <application>` and application supporting components ({ref}`machines <machine>`, {ref}`storage <storage>`, {ref}`network spaces <space>`, {ref}`relations <relation>`, and so on).

A model is created by a {ref}`user <user>`, and owned in perpetuity by that user (or a new user with the same name), though it may also be used by any other user with model access level, within the limits of their level.

A model is created on a {ref}`controller <controller>`.  Both the model and the controller are associated with a {ref}`cloud <cloud>` (and a cloud {ref}`credential <credential>`), though they do not both have to be on the same cloud (this is a scenario where you have a 'multicloud controller'). Any entities added to the model will use resources from that cloud.

One can deploy multiple applications to the same model. Thus, models allow the logical grouping of applications and infrastructure that work together to deliver a service or product.  Moreover, one can apply common {ref}`configurations <configuration>` to a whole model. As such, models allow the low-level storage, compute, network and software components to be reasoned about as a single entity as well.

A controller does the work for each model with a set of {ref}`model workers <model-workers>`, which it starts when the model is activated.

(model-workers)=
## Model workers

On machine clouds a model has no process of its own: the {ref}`controller agent <controller-agent>`'s model worker manager starts an engine of **model workers** for every activated model, and those workers do the controller's work for the model. They request machines from the model's cloud, remove what has been removed, move the model to another controller and keep its records tidy. On Kubernetes the model also gets a model operator, a pod in the model's namespace that represents the model inside the cluster.

In a controller with several {ref}`high-availability <high-availability>` nodes, every node starts the engine and one node at a time runs the workers: the node that holds the model's responsibility lease. The workers pause while the model migrates and stop when it is dead. Those that call the cloud run only while the model's {ref}`credential <credential>` is valid.

- **Machine provisioning (`compute-provisioner`, `instance-poller`, `firewaller`):** On machine models, request machines from the cloud and record the instance ID and addresses, poll the cloud for instance state, and manage the model's firewall rules (see {ref}`machine provisioning <the-machines-machinery>`).
- **Kubernetes provisioning (`caas-application-provisioner`, `caas-firewaller`, `caas-model-config-manager`, `caas-model-operator`):** On Kubernetes models, create, manage and destroy the Kubernetes resources of each application to match its requested state, and deploy the model operator.
- **Storage (`storage-provisioner`):** Provisions, attaches and detaches the model-scoped volumes and filesystems of the model's {ref}`storage <storage>`. Each machine agent runs one for machine-scoped storage.
- **Removal (`removal`):** Runs the scheduled removal jobs that delete the records of removed entities (see {ref}`model removal <the-model-removal>`).
- **Migration (`migration-master`):** Drives the model's migration to another controller when the model is the source (see {ref}`model migration <the-model-migration>`).
- **Charms (`charm-revisioner` and `async-charm-downloader`):** Poll Charmhub every 24 hours for new revisions of the model's charms, so that status can show an update, and download charms that are deployed but not yet downloaded.
- **Cross-model relations (`remote-relation-consumer`):** Manages the relations, and their settings, where one end is a remote application.
- **Secrets (`secrets-pruner` and `user-secrets-drain-worker`):** Prune obsolete user secret revisions, and move user secrets from the inactive secret backend to the active one.
- **Housekeeping (`operation-pruner`, `change-stream-pruner`, `unitless`):** Prune stored operation results according to the model configuration, prune the model's change stream, and run the scriptlets of {ref}`applications <application>` that have no units.

(the-models-persistence)=
(the-model-record)=
(the-model-persistence-rules)=
(the-types-of-model)=
(types-of-model)=
(iaas-caas-models)=
(the-controller-model)=
(regular-model)=
## Model data model

A model has records in two places. The controller database holds the model's identity and its life, and the model's own database holds a copy of that identity with everything the model contains.

- **Controller record (`model` in the controller database):** The authoritative record lives in the controller database, identified by its natural key, the model's name plus the {ref}`user <user>` that qualifies it (the **qualifier**, which disambiguates same-named models of different owners). The record carries the model's type, its cloud, region and credential, and the controller-model flag (see {ref}`types of model <types-of-model>` for what each of those discriminators does).
  - *Rule:* The natural key is unique: a model with the same name and owner cannot be added twice.
  - *Related errors:*
    - **`model already exists`**: *Trigger:* Adding a model whose name and owner, or whose UUID, match an existing model. *Remediation:* Choose another name or owner.
    - **`model not found`**: *Trigger:* Querying a model by UUID or name that does not exist. *Remediation:* Verify the model exists on the controller.
    - **`secret backend already set`**: *Trigger:* Setting a secret backend on a model record that already has one. *Remediation:* Update the backend instead of setting it again.
    - **`user not found on model`**: *Trigger:* Looking up a user's grant on a model when none exists. *Remediation:* Grant the user access on the model.
- **Life and activation:** The model is created alive. Its life is carried in the **controller** database, the shared alive / dying / dead cycle every entity has, plus an *activated* flag that says the model is live on its controller: a model imported by a {ref}`migration <the-model-migration>` starts unactivated and is activated only once its agents have re-oriented. The model's removal moves it to dying and dead; the record stores it. There is no life on the model's own record in the model database: life belongs to the controller database.
  - *Related errors:*
    - **`model not activated`**: *Trigger:* Operating on an imported model before its agents have reported success. *Remediation:* Wait for the migration's activation.
    - **`model already activated`**: *Trigger:* Activating a model that is already live on its controller. *Remediation:* None; the model is already serving.
- **Namespace record:** Names the model's own Dqlite database.
  - *Related error:*
    - **`model namespace not found`**: *Trigger:* The model's Dqlite namespace record does not exist. *Remediation:* Verify the model's database exists.
- **Model database copy:** The model's stored footprint is thin, and deliberately so: everything the model *contains* lives in the entities' own records (see {ref}`the database <database>`). The model's own **model database** keeps a read-only, denormalized copy of its identity (its UUID, its controller, its name and owner, its type, its cloud, region and credential, and the controller-model flag), one record per model database, enforced by the schema, and its model-life record mirrors the controller-side life as a best-effort facsimile; it is what the model-side processes read.
- **Model's own records:** The model database also carries the model's own records: its {ref}`configuration <model-configuration>`, its {ref}`constraints <constraint>`, its {ref}`storage pools <storage>` and the target agent version (the agent version the model's machines should run and the latest one known; see {ref}`upgrading your deployment <upgrade-your-deployment>`).
  - *Related error:*
    - **`model constraints not found`**: *Trigger:* Querying the model's constraints record when it does not exist. *Remediation:* Check that the model has constraints set.
- **Type (a column of `model`):** Derived from the cloud at creation: a Kubernetes cloud makes the model `caas`, and every other cloud makes it `iaas`. A model is one or the other, never both. The type selects the model's workers, which {ref}`secret <secret>` backend its secrets live in, and whether its applications scale by {ref}`pods <application>` or {ref}`machines <machine>`.
- **Controller-model flag (a column of `model`):** A role, not a type. It marks the model that runs Juju itself.
  - **The controller model (`controller`):** This is your Juju management model. A Juju deployment will have just one controller model, which is created by default when you create a controller (`juju bootstrap`). It typically contains a single machine, for the controller (the `controller` application). If controller {ref}`high availability <high-availability>` is enabled, then the controller model would contain multiple instances. The `controller` model may also contain certain applications which it makes sense to deploy near the controller, for example the `juju-dashboard` application.
  - **Regular model:** This is your Juju workload model. A Juju deployment may have many different workload models, which you create explicitly. It is the model where you typically deploy your applications.

**Writers:** The model service in the controller performs the writes: creating a model writes its records in both databases; configuring it rewrites the configuration records; migrating it creates an *importing* model, unactivated until the migration's agents report success (see {ref}`Model migration <the-model-migration>`).

(the-model-watchers)=
## Model change streams

Juju does not poll for changes to models. Its services publish change streams over the model records, and the workers that act on models subscribe to them.

The model domain's watchable service exposes these watch surfaces:

- **All models:** The controller database's model records: the
  Undertaker's surface for models becoming dead.
- **Activated models:** The models that have finished activating on
  the controller: it drives the model worker manager, which starts an
  engine of {ref}`model workers <model-workers>` for each.
- **Model migration deletions:** The Undertaker's second surface:
  models deleted by a migration's reap phase.
- **One model:** Changes to a single model's record.
- **One model's cloud credential:** The credential a model tracks,
  so the model can react when the user changes it.

Every watcher fires once immediately when it is created, the initial
query being the baseline snapshot, and again on each qualifying change
(see {ref}`the watcher pattern <watchers>`).

(the-models-declaration)=
(the-model-declaration-rules)=
(the-models-execution)=
(the-model-execution-rules)=
## Model operations

A client asks the controller to create, configure, migrate or remove a model. The controller records the request, and the model's {ref}`workers <model-workers>`, which watch the model's records, carry out the long-running parts.

(model-creation)=
### Model creation

Creating a model requires {ref}`cloud add-model access <user-access-cloud-add-model>`. By the time the call returns, the records exist and the model's workers are running.

- *Rule:* The model name is not empty, and the qualifier (its owner) must be a valid user name: `admin`, or `alice@external` for a {ref}`user <user>` from an identity provider.
- *Rule:* A model whose name and owner match an existing model is refused with `model already exists` (see {ref}`the model record <the-model-record>`).
- *Related errors:*

  ```{list-table}
  :header-rows: 1
  :widths: 30 40 30

  * - Message
    - Trigger
    - Remediation
  * - `model credential not valid`
    - The cloud credential for the model does not validate.
    - Check the credential against the model's cloud.
  * - `agent version not supported`
    - The model's target agent version is not supported.
    - Target a supported agent version.
  * - `agent stream not valid`
    - The model's agent stream is not valid.
    - Use a valid agent stream.
  ```

### Model configuration

A client sets model configuration through the controller API.

- *Rule:* Model configuration keys must be known keys, and their values must pass the schema validation (a `validation error` naming the invalid attributes).

```{ibnote}
See also: {ref}`Juju | Manage models <manage-models>`, {ref}`Juju | Configure a model <configure-a-model>`, {ref}`Terraform Provider for Juju | Manage models <tfjuju:manage-models>`
```

(the-model-migration)=
### Model migration

A model can be moved between controllers with the migrate operation
(for example, `juju migrate <model> <target-controller>`). The
migration is a phase machine driven by the source controller: the
source's {ref}`migration master <model-workers>` worker quiesces the model's agents (locks them down), sends
the model envelope (the YAML export of the model database plus the
model's controller-DB data) to the target controller, whose import
runs as ordered operations with rollback, then lets the agents
validate against the target and rewrite their agent configuration to
re-orient to it. Once the agents report success, the source
activates the imported model on the target, transfers the model's
logs, and reaps the source model (active users are redirected). An
aborted migration rolls the phases back on both sides.

- *Rule:* Importing a model from a 3.6 controller transforms what 4.1 stores differently: a 3.6 consumer proxy carries one offer connection per relation, and the import creates a synthetic application, an offer connection and a consumer record for each of them; annotations and the controller's API addresses land in their 4.1 tables; and secret grants whose relation key no longer resolves are skipped with a warning.
- *Rule:* The log transfer checkpoints its progress in the migrated model's database. An interrupted transfer resumes from the last checkpoint and may replay records after it, so duplicates are possible and no log is missed.
- **Model migration:** Migrating moves a model to another controller (for example, with `juju migrate`). The
  source controller's migration master walks the stages below. At each step it
  writes a status message, and the controller reports that message, with the
  migration's start and end times, as part of the model's information. A client
  shows it; for example, `juju show-model` prints it in its `migration` field. If any stage fails, the
  migration aborts and the model stays on the source controller.

  The stages, in order, with the status messages each one writes:

  ```{list-table}
  :header-rows: 1
  :widths: 16 44 40

  * - Stage
    - Status message shown
    - What happens
  * - Prechecks
    - `performing source prechecks`, `performing target prechecks`, `quiescing, waiting for agents to report back`
    - Both controllers check that the migration can proceed. The model's agents then stop their work and are barred from restarting it (the migration code calls this *quiescing*), and the checks repeat once the model is stable.
  * - Transfer
    - `exporting model`, `importing model into target controller`, `uploading model binaries into target controller`
    - The source exports the model; the target imports it and receives its binaries.
  * - Validation
    - `validating, waiting for agents to report back`, `checking machines in migrated model`, `activating model in target controller`
    - The agents validate against the target, the machines are compared with the cloud, and the target activates the model.
  * - Handover
    - `successful, waiting for agents to report back`, `transferring ownership of cloud resources to target controller`
    - The agents re-point to the target and the target takes over the cloud resources.
  * - Log transfer
    - `successful, transferring logs to target controller (<n> sent)`, then `successful, transferred logs to target controller (<n> sent)`
    - The model's logs are copied to the target.
  * - Cleanup
    - `successful, removing model from source controller`
    - The source copy is removed.
  * - Rollback
    - `aborted, removing model from target controller: <reason>`
    - A failure in any earlier stage lands here: the target copy is removed and the model stays on the source.
  ```

  In the migration code these stages are the phases `QUIESCE`, `IMPORT`,
  `VALIDATION`, `SUCCESS`, `LOGTRANSFER`, `REAP` and `ABORT`. They are internal:
  the client reports the status messages, not the phase names.

  When a stage fails, the status message becomes the final `aborted, ...: <reason>`
  line, and the `<reason>` is one of the messages in the table below. The Stage
  column names the stage of the table above that raises it, and the rows follow
  the order of the stages. A reason is often prefixed by the step that failed
  (`model data transfer failed, <inner>`, `model activation failed, <inner>`).

  ```{list-table}
  :header-rows: 1
  :widths: 14 30 30 26

    * - Stage
      - Message as seen
      - Trigger
      - Remediation
    * - Prechecks (source)
      - `model is <life>`
      - The model is dying or dead.
      - Wait for it to be removed, or cancel its removal.
    * - Prechecks (source)
      - `model is being imported as part of another migration`
      - The model is the target of an unfinished migration.
      - Wait for that migration to finish or abort.
    * - Prechecks (source)
      - `model has revoked credentials`
      - The model's cloud credential is revoked.
      - Set a valid credential on the model, then retry.
    * - Prechecks (source)
      - `machine <id> is <life>`
      - A machine (a model machine on the source side, a controller machine on either side) is not alive.
      - Let the removal finish or bring a machine back, then retry.
    * - Prechecks (source)
      - `unit <unit> hasn't joined relation "<key>" yet`
      - A unit has not entered scope in one of the model's relations.
      - Wait until the unit has joined, then retry.
    * - Prechecks (source)
      - `all charms now require a manifest.yaml file, this model hosts charm(s) with no manifest.yaml file: <list>`
      - The model runs a charm revision without a manifest.
      - Refresh the listed applications to a revision that has one.
    * - Prechecks (source)
      - `fan networking not supported, remove fan-config "<value>" from migrating model config`
      - The model config sets `fan-config` or `container-networking-method` (fan networking).
      - Remove the key from the model config.
    * - Prechecks (source)
      - `cannot migrate to controller due to issues:` followed by a list
      - The upgrade validators found blockers for this model.
      - Resolve each listed issue.
    * - Prechecks (source)
      - `controller: upgrade in progress`
      - The source controller is upgrading (the controller checks on the source side carry the `controller:` prefix).
      - Wait for the upgrade to finish.
    * - Prechecks (target)
      - `model has higher version than target controller (<model> > <controller>)`
      - The model's agent version is above the target controller's.
      - Upgrade the target controller first.
    * - Prechecks (target)
      - `source controller has higher version than target controller (<src> > <tgt>)`
      - The controllers' major or minor versions are incompatible.
      - Use a target controller at the same or a newer version.
    * - Prechecks (target)
      - `model's cloud "<cloud>" not found on target controller`
      - The target does not have the model's cloud.
      - Add the cloud to the target controller.
    * - Prechecks (target)
      - `model named "<name>" already exists` / `model with same UUID already exists (<uuid>)` / `model is being migrated out of target controller`
      - The target holds a model with the same name or UUID, or one that is being migrated away.
      - Rename or remove the clashing model; wait if one is leaving.
    * - Prechecks (target)
      - `target controller does not support the model migration format; upgrade the target controller to Juju 4.1 or later`
      - The target's migration-target facade is older than version 8.
      - Upgrade the target controller to 4.1 or later.
    * - Validation
      - `machine sanity check failed, <n> error(s) found`
      - After the transfer, the machines Juju lists differ from what the cloud reports.
      - Reconcile the cloud and Juju (for example remove the leftover instance), then retry.
    * - Prechecks, Validation, Handover
      - `<step>, timed out waiting for agents to report` / `<step>, some agents reported failure`
      - Agents did not report in time, or reported an error. `<step>` is `quiescing`, `validating` or `successful`.
      - Raise `migration-agent-wait-time`; read the agents' logs.
    * - Cleanup
      - `removing exported model failed: <inner>`
      - The source copy could not be removed after the model went live on the target.
      - Remove the leftover source model; the migrated model is already live.
  ```

  Errors that come from the target's import (a missing owner, an unknown
    credential) are wrapped in `model data transfer failed, <inner>`; their inner
    messages are listed on the user and credential pages (`user not found`,
    `credential not found`) and linked from here.

- *Rule:* The target must be a Juju 4.1 (or newer) controller: it must support the migration envelope facade version the source sends (version 8 of the migration target facade). A 4.0 controller cannot be a migration target.
- *Related errors:*
  - **`target controller does not support the model migration format; upgrade the target controller to Juju 4.1 or later`**: *Trigger:* The prechecks find a target whose migration target facade is older than the envelope needs. *Remediation:* Upgrade the target controller to Juju 4.1 or later.
  - **`model not redirected`**: *Trigger:* Looking up a model's completed redirect after a migration when none exists. *Remediation:* Verify the migration finished its redirect phase.

(the-model-removal)=
### Model removal

A model is destroyed with `juju destroy-model`. Destruction is the
largest removal and runs the same cooperative pattern at scale: the
controller marks the model Dying and wakes the {ref}`Undertaker <controller-agent>` worker,
which tears the model down in dependency order (applications first,
then the machines, and the model's Dqlite database last, as the final
act after all cloud resources are released). The controller never
deletes its own database; the Undertaker does.

- *Rule:* The controller never deletes its own database: the model's Dqlite database goes last, after all cloud resources are released.
- *Rule:* The cascade schedules a removal job for every machine of the model, container hosts included. A host's job reports itself incomplete and is polled again until its units and child machines are gone, then deletes the host; the model's job cannot finish while any machine record remains.
- *Rule:* The controller model is removed only with the force option, and a model that holds persistent storage is removed only when the request says what to do with the storage.
- *Related errors:*
  - **`force required for removal job`**: *Trigger:* Removing the controller model, or an offer that still has connections, without the force option. *Remediation:* Force the removal.
  - **`model has an active migration import claim`**: *Trigger:* A generic removal targets a model that is mid-import. *Remediation:* Abort the import; it tears the imported model down.
  - **`model has persistent storage`**: *Trigger:* Removing a model that holds persistent storage without stating whether to destroy or release it. *Remediation:* Pass the option that destroys or releases the storage.
