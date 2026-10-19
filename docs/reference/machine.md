---
myst:
  html_meta:
    description: "Juju machine reference: compute resources on bare metal, VMs, and LXD containers. The records Juju stores about a machine, the change streams it publishes, and the operations that request, provision, reprovision, remove and reboot machines."
---

(machine)=
# Machine

```{ibnote}
See also: {ref}`manage-machines`
```

```{ggarch}
:file: ../machine.ggarch
:view: Machine and its neighbours
:caption: A machine belongs to one model, which has many machines. A machine hosts one or more units, and each unit belongs to one application. A container machine runs inside one host machine.
:alt: A model has many machines. A machine hosts one or more units, and each unit belongs to one application. A container machine runs inside one host machine.
```

In Juju, a **machine** is a {ref}`compute resource <resource-compute>` from a machine {ref}`cloud <cloud>`: a bare cloud instance, or a {ref}`LXD container <cloud-lxd>` on one, requested implicitly by a {ref}`deploy <deploy-an-application>` or {ref}`add-unit <command-juju-add-unit>` {ref}`placement <placement-directive>`, or explicitly through the {ref}`controller <controller>` API.

```{important}

A LXD container is named after its host machine: `0/lxd/5` is LXD container `5` on machine `0`. Containers appear under 'Machines' in status outputs.

```

A machine belongs to one {ref}`model <model>` and hosts the {ref}`units <unit>` of one or more {ref}`applications <application>`, usually one. A container machine runs inside a host machine, and a host can run many containers. A machine's cloud instance sits on the model's cloud.

(the-machine-in-the-data-model)=
(machine-persistence-rules)=
## Machine data model

```{ggarch}
:file: ../machine.ggarch
:view: Machine data model
:caption: A machine is one primary record that holds its identity, its agent's credentials and its keep-instance flag, and it points at one net node. At most one of each of these records points back at it: status, constraint, cloud instance (with at most one status of its own), manual flag, parent and reprovision request. A host machine can be the parent of many containers, and a cloud instance may name an availability zone that many instances share. The platform record also points at the machine. Brackets list each record's descriptive columns, and every line is a foreign key in the 4.1 model schema.
:alt: The machine record points at one net node. The platform, constraint, status, cloud instance, parent, manual and reprovision request records point at the machine, at most one of each except the platform record, and a host machine can be the parent of many containers. The cloud instance record has at most one status record and may point at an availability zone.
```

Juju stores everything it knows about a machine in the model's {ref}`database <database>`. One primary record holds the machine's identity, and a few smaller records hold the rest; each of them points back at the primary record.

- **Machine (`machine`):** The primary record. It holds an internal ID that the other records point at, the name the client types (`0`, `1/lxd/0`), the agent's credentials (a nonce and a password hash), the hostname, and the keep-instance flag, which says whether the cloud instance survives the machine's removal.
  - *Rule:* A machine's name is unique per model. Machine numbers are `0` or positive integers without leading zeros, and containers nest one level deep (`0/lxd/4`, never `0/lxd/0/lxd/0`).
  - *Related errors:*
    - **`machine name "<name>" has too many containers`**: *Trigger:* A machine name with more than one container level. *Remediation:* Name the container directly under its host machine.
    - **`machine not found`**: *Trigger:* A lookup by machine UUID or name that matches no machine. *Remediation:* Check the model's machines.
    - **`machine already exists`**: *Trigger:* A model migration imports a machine whose name the target model already holds. *Remediation:* Migrate into a model that has no machine of that name.
- **Life (a column of `machine`):** `alive`, `dying` or `dead`. The cloud instance record has a life of its own, which gates deleting the machine.
  - *Rule:* A machine becomes `dead` only when it is `dying` and hosts no containers, units or storage.
  - *Related error:*
    - **`machine is dead`**: *Trigger:* An operation on a machine whose life is `dead`. *Remediation:* None; the machine's lifecycle has ended.
- **Net node (`net_node`):** The network identity that a machine shares with every unit running on it.
- **Platform (`machine_platform`):** The machine's base: its operating system and channel, and its architecture.
- **Constraint (`machine_constraint`):** Points the machine at its {ref}`constraint <constraint>` record.
- **Manual flag (`machine_manual`):** Exists only for a machine the user provisioned by hand over SSH.
- **Parent (`machine_parent`):** Exists only for a container. It holds two pointers, the container and its host, and the container's type is recorded alongside. Removal, model migration and recovery also read it.
  - *Related error:*
    - **`machine has no parent`**: *Trigger:* Asking for the parent of a machine that is not a container. *Remediation:* Check that the machine is a container (`<host>/lxd/<n>`).
- **Cloud instance (`machine_cloud_instance`) and zone (`availability_zone`):** The record of the instance behind the machine: its creation state, its hardware (arch, CPU, memory, root disk, virtualisation type) and its availability zone. It is created with the machine, and its instance ID stays empty until the cloud reports one. Removal, status, storage provisioning, network, model migration, recovery and agent passwords also read it.
  - *Related errors:*
    - **`machine not provisioned`**: *Trigger:* An action that needs the instance runs before the cloud has reported an instance ID. *Remediation:* Wait until the compute provisioner has recorded the instance.
    - **`machine cloud instance already exists`**: *Trigger:* Recording an instance for a machine that already has one. *Remediation:* Check for a second provisioner acting on the same machine.
- **Status (`machine_status`) and instance status (`machine_cloud_instance_status`):** Two status records.
  - The agent's health report, written by the machine agent: `pending`, `started`, `stopped`, `error`. `down` is never written, only read when the agent has not been seen recently.
  - The provisioning lifecycle, written by the controller and by the provider or container broker: `pending`, `allocating`, `running`, `provisioning error`, `unknown`. `provisioning error` is a status value: cloud and container providers record it when allocation fails.
- **Reprovision request (`machine_reprovision`):** Exists while a reprovision request is pending: the machine's name and the time of the request.
  - *Rule:* A machine has at most one pending reprovision request.
  - *Related error:*
    - **`reprovision request already exists`**: *Trigger:* A reprovision request is already pending for the machine. *Remediation:* Wait for the pending request to complete.

Machine kinds are not stored. Juju works them out from the records: a container has a parent record, a manual machine has a manual flag, and a controller machine hosts a unit of the `controller` application, runs the {ref}`controller agent <controller-agent>` and runs every model worker.

The machine service inserts the records (`AddMachine`, `SetMachineCloudInstance`), and the removal service carries out the teardown.

(machine-watchers)=
## Machine change streams

Juju does not poll for changes to machines. Its services publish change streams over the machine records, and the workers and agents that act on machines subscribe to them.

- **Machine life and dependants:** Drives the local shutdown watcher.
- **Container life:** Drives host-side container provisioning loops.
- **Model machines:** Drives the compute provisioner.
- **Instance states:** Drives the instance poller's steady-state mirroring (life, start times, cloud instances) and the reboot machinery.

Every watcher fires an initial baseline snapshot on creation, followed by notifications on qualifying changes. See {ref}`the watcher pattern <watchers>`.

(the-machines-machinery)=
(machine-execution-rules)=
(machines-and-units)=
## Machine operations

Machine operations are shared between the controller and the host machine. The controller runs background workers that request, provision and remove machines, and the machine agent on each host looks after its own lifecycle and the workload placed on it.

(machine-declaration-rules)=
### Machine requests

A client asks the controller for compute in two ways. An **explicit** request creates machines through the controller API, gated on model {ref}`write access <user-access-model-write>`. An **implicit** request is a deploy or add-unit placement, and most machines arise this way. A machine usually hosts one unit, and placing several units of the same or of different applications on one machine saves resources.

- **Target designations:** Comma-separated lists naming new or existing targets.
  - `0`: an existing machine; `0,4`: machines 0 and 4.
  - `lxd`: a new container on a new host; `lxd:25`: on machine 25.
  - `0/lxd/4`: the existing container `4` on machine `0`.
  - An empty designation: a new machine.
  - `lxd` or `lxd:25` with `virt-type=virtual-machine`: a VM instead of a container.
  - `3,0/lxd/2,lxd:5`: machine 3, container 2 on machine 0, and a new container on machine 5.
  - Machine clouds only; rejected on Kubernetes models.
  - *Rule:* Containers nest one level deep, so a name with more than one container level is refused with `machine name "<name>" has too many containers`.
  - *Rule:* The only container type in 4.1 is `lxd`.
  - *Related error:*
    - **`invalid container type "<type>"`**: *Trigger:* A designation or a container constraint names a container type other than `lxd`. *Remediation:* Use `lxd`.
- **Base:** A machine's **base** identifies its operating system image: the OS name, `@`, and the OS channel as `<track>` or `<track>/<risk>` (`ubuntu@22.04`, `ubuntu@22.04/stable`). A base replaces the older notion of a series.
- **Machine constraints:** Constraints can be set directly on a machine, as well as on the model or an application. Only an architecture different from the recorded instance's is compared.
  - *Related errors:*
    - **`invalid machine constraints`**: *Trigger:* The constraints name a container type or a space that does not exist. *Remediation:* Correct the container type or the space.
    - **`machine constraint violation`**: *Trigger:* The constraints are set on a machine whose recorded cloud instance has a different architecture from the one required; architecture is the only characteristic compared. *Remediation:* Require the architecture the machine has, or request a new machine.

```{ibnote}
See also: {ref}`Terraform Provider for Juju | Manage machines <tfjuju:manage-machines>`
```

### Machine provisioning

Provisioning turns a machine record into a running machine. It happens in one of three ways:

- **Controller-side provisioning:** The compute provisioner, driven by the model-machines {ref}`change stream <machine-watchers>`, monitors unprovisioned machines, requests cloud instances, and records the instance ID and addresses as they become known.
- **Agent-side provisioning:** Host machine agents provision and manage local LXD containers through the LXD broker, driven by the container-life {ref}`change stream <machine-watchers>` over the controller API.
- **Manual provisioning:** The controller renders a setup script, executed externally over SSH, after which the host agent registers itself.

### Machine reboot

The reboot machinery, driven by the instance-states {ref}`change stream <machine-watchers>`, decides whether a machine shuts down, reboots or does nothing: a machine whose parent needs a reboot shuts down, a machine that needs a reboot itself reboots, and otherwise nothing happens.

- **Nesting:** The reboot lookup supports one level of container nesting.
  - *Related error:*
    - **`grandparent machine are not supported currently`**: *Trigger:* The lookup finds a machine whose parent itself has a parent. *Remediation:* None; name validation refuses containers nested deeper than one level, so the check guards a layout that Juju does not support.

### Machine reprovisioning

Reprovisioning replaces the cloud instance behind a machine whose instance is lost, and keeps the machine itself. A client asks the controller for it (for example, with `juju reprovision-machine`), and the controller works in three stages.

- **Eligibility:** Only an `alive` machine that is not a container, a controller machine or a manually provisioned machine, hosts no containers, has no model-scoped storage attached, and has no reprovision request pending (`reprovision request already exists`) is eligible.
  - *Related errors:*
    - **`machine is not alive`**, **`machine is a container`**, **`machine is a controller machine`**, **`machine is manually provisioned`**, **`machine hosts child containers`**: *Trigger:* The machine is `dying` or `dead`, or is of a kind that cannot be reprovisioned, or hosts containers. *Remediation:* Reprovision only an `alive` machine that is a plain, provider-provisioned host.
    - **`model-scoped storage attached to machine`**, **`storage scope is ambiguous`**, **`machine storage is not alive`**: *Trigger:* The machine's attached storage cannot be reset safely. *Remediation:* Detach or remove the storage first.
- **Liveness checks:** The machine agent must be absent, and the provider must report the instance missing or not running.
  - *Related errors:*
    - **`machine agent is still present`**: *Trigger:* The machine agent still reports presence. *Remediation:* Stop the agent or confirm the machine is lost.
    - **`machine provider instance is running`**: *Trigger:* The provider reports the instance as running. *Remediation:* Stop or remove the instance at the provider, then retry.
- **Replacement:** In one transaction, the controller detaches the lost instance; gives the machine's units new ordinals, so the replacement unit agents start with fresh identities; departs the units from their relation scopes and archives their settings, so the remaining units see a consistent departure; resets the machine-scoped storage state; and returns the machine and its instance to `pending`. The compute provisioner then provisions the replacement instance.
  - *Related error:*
    - **`machine cloud instance changed`**: *Trigger:* The machine's instance differs from the one checked, because a replacement appeared in between. *Remediation:* Retry the request.

### Machine removal

```{ggarch}
:file: ../machine.ggarch
:view: Machine life
:caption: A machine's life during removal. A machine is alive until its removal is requested, and then dying. It becomes dead once it hosts no containers, units or storage. The removal job deletes its records when its cloud instance is dead. A forced removal deletes the records from dying, without waiting.
:alt: A machine is alive, then dying when removal is requested, then dead once it hosts no containers, units or storage, then removed when the removal job deletes its records and the cloud instance is dead. A forced removal goes from dying straight to removed.
```

Removing a machine is cooperative:

- **Request:** The controller marks the machine `dying`.
- **Shutdown:** The machine's agent learns of it through the machine-life {ref}`change stream <machine-watchers>` and shuts the machine down once its units and storage are gone.
- **Marking dead:** Once the machine hosts no containers, units or storage, it is marked `dead`.
  - *Related errors:*
    - **`entity still alive`**: *Trigger:* Marking a machine `dead` while its life is still `alive`. *Remediation:* Start the removal first, so the machine becomes `dying`.
    - **`machine has containers`**, **`machine has units`**, **`machine has storage`**: *Trigger:* Marking a machine `dead` while it still hosts containers, units or storage. *Remediation:* Remove or move them first.
- **Deletion:** Scheduled removal jobs purge the records, and the keep-instance flag decides whether the cloud instance is released.
