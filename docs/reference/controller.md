---
myst:
  html_meta:
    description: "Juju controller reference: the control plane that serves many models, on one or several clouds. The controller data model (the controller record, models, high-availability nodes and configuration, and the client's own store), and the operations that bootstrap, configure, expand and remove it."
---

(controller)=
# Controller

```{ggarch}
:file: ../controller.ggarch
:view: Controller and its neighbours
:caption: A controller serves many models. Each model is on one cloud, which has to be a cloud the controller knows, and uses one credential for that cloud where the cloud needs one. The controller itself runs on one cloud and knows many.
:alt: A controller serves many models. A model is on one cloud and uses one credential, where one is needed. A credential is for one cloud. The controller runs on one cloud and knows many clouds.
```

In Juju, the **controller** is the control plane: the running system that implements every change a {ref}`user <user>` asks for.

A controller comes into being through {ref}`bootstrap <bootstrap-a-controller>`, which provisions a machine on a {ref}`cloud <cloud>`, installs the {ref}`controller agent <controller-agent>` (a {ref}`jujuagentd <jujuagentd>` process whose workers run the API server and an embedded [Dqlite](https://canonical.com/dqlite) database) and starts it. That cloud is the controller's own.

Once it is running, clients and agents reach the controller through its **API server**. Clients (the `juju` CLI, the Terraform provider and other tools) and agents all connect to it, and it authenticates each request and routes it to the versioned endpoints for that kind of caller.

From then on the controller creates and manages {ref}`models <model>`. Each model is on one cloud, and the controller requests the model's resources from that cloud. The cloud can be the controller's own or another one, but it has to be a cloud the controller knows, with a {ref}`credential <credential>` where the cloud needs one. So a single controller can manage models on several clouds, and one cloud can host many models.

The {ref}`controller model <the-controller-model>` is the one model that runs Juju itself. The controller knows its users, clouds and credentials as a whole, while each model's {ref}`applications <application>`, {ref}`machines <machine>` and {ref}`units <unit>` are stored with that model.

(the-controller-persistence-rules)=
## Controller data model

```{ggarch}
:file: ../controller.ggarch
:view: Controller data model
:caption: A controller's state in the controller database is one controller record, the models it serves and its high-availability nodes. A controller serves many models; each model is on one cloud and a cloud can host many models, so one controller can serve models on several clouds. Each model has one namespace record naming the separate Dqlite database that holds its records. Each node has at most one agent version and one password, and many addresses of each kind: API, agent, client and peer. The object store backend has at most one S3 configuration. The controller record names the controller model by UUID, without a database foreign key, and the controller configuration is a standalone set of key/value records. Brackets list each record's descriptive columns, and every other line is a foreign key in the 4.1 controller schema.
:alt: The controller record points at the controller model. Each model points at one cloud, and many models can share a cloud. A model is pointed at by one namespace record. A controller node is pointed at by at most one agent version record and one password record, and by many API, agent, client and peer address records. An object store backend is pointed at by at most one S3 configuration record.
```

A controller's state is stored in two places: first on the machine where you run the Juju client, and then in the controller's own {ref}`database <database>`. The client machine keeps a small set of files that record the controllers, accounts, models, clouds and credentials the client knows about, so that the client can reach a controller once it exists. The controller database holds the controller itself, the models it serves, its high-availability nodes and its configuration.

### Controller client store

The client keeps these files in its data directory:

- `clouds.yaml`: the clouds the client knows.
- `credentials.yaml`: the credentials it uses on them.
- `controllers.yaml`: the controllers it has bootstrapped or logged in to.
- `accounts.yaml`: the accounts it uses on those controllers.
- `models.yaml`: the models it has seen.

Bootstrap writes the new controller's entry into these files, and later client commands read it to find the controller.

### Controller database records

The controller itself is stored as one record, and the records of everything it manages hang off it:

- **Controller (`controller`):** The controller itself. It holds the controller's UUID, the controller model it lives in, the agent version an {ref}`upgrade <upgrade-your-deployment>` targets, the API port, and the controller's TLS certificate and system identity.
  - *Rule:* There is exactly one controller record, and the database refuses a second. Adding high-availability {ref}`nodes <high-availability>` adds nodes to the one controller and never creates another controller.
  - *Rule:* The controller has no life or status of its own. The liveness of its nodes is the Dqlite cluster's business.
  - *Related error:*
    - **`controller agent info not found`**: *Trigger:* The controller agent's recorded information is missing. *Remediation:* Check that the controller has finished bootstrapping.
- **Model registry (`model` and `model_namespace`):** One record per model the controller serves. Each is paired with the name of the separate Dqlite database where that model's records live, and the controller record points at the controller model, the one model that runs Juju itself. Juju works out the controller model, the controller cloud and the controller {ref}`machine <the-machine-in-the-data-model>` from the records around them instead of flagging them.
- **The `controller` application:** The controller is also an {ref}`application <application>`: the one named `controller`, deployed in the controller model from the `juju-controller` {ref}`charm <charm>`. A single-record marker identifies it, and the database allows only one such marker.
  - *Charm origin:* The charm archive (`controller.charm`) ships in the agent's data directory, and bootstrap deploys it from a local origin at revision 0. If the archive is missing, bootstrap falls back to Charmhub.
  - *Units:* The application runs real units: one created at bootstrap, and one per {ref}`HA node <high-availability>`. A machine is the controller machine when it hosts a unit of this application.
- **HA nodes (`controller_node`):** One record per Dqlite node in the {ref}`high-availability <high-availability>` cluster, with the node's Dqlite identity and bind address, each unique, and a life of its own (`alive`, `dying` or `dead`). Further records per node hold the agent version it reports (`controller_node_agent_version`), its password (`controller_node_password`), and the addresses it serves on: the API addresses (`controller_api_address`), which let agents re-orient as the cluster changes, and separate addresses for agents, clients and peer controllers (`controller_agent_address`, `controller_client_address`, `controller_peer_address`). High availability does not create a second controller; it adds nodes that run the same controller's database and API. Agent passwords, upgrades, model agents and the change stream also read the node records.
  - *Related errors:*
    - **`controller not found`**: *Trigger:* A lookup of a controller node, for example by an agent authenticating or by the upgrader, names a node that does not exist. *Remediation:* Check the node's identifier against the controller's nodes.
    - **`controller address not valid`**: *Trigger:* A controller address is recorded as empty. *Remediation:* Provide a non-empty address.
    - **`no controller IDs found`**: *Trigger:* A query for controller nodes finds none. *Remediation:* Check that the controller's nodes are recorded.
    - **`no API addresses found`**: *Trigger:* The high-availability request finds no API addresses for the controller nodes. *Remediation:* Check that the nodes have reported the addresses they serve on.
    - **`stale controller membership`**: *Trigger:* Address selection used a controller membership that changed after it was read. *Remediation:* Retry the request.
- **Configuration (`controller_config`):** A key/value set of controller-wide settings. Reading it returns the stored keys together with values that Juju computes from the controller record. Provisioning, key updating and model migration also read it.
- **Storage backends (`object_store_backend` and `object_store_backend_s3_config`):** The controller's persistent data lives in the Dqlite databases and in blob storage. Both default to the controller's filesystem (the `file` object-store backend). An S3-compatible object store such as AWS S3, MicroCeph, or MinIO can take over blob storage, and then the controller database keeps the object-store backend records, with the backend type and, for S3, its endpoint, credentials and region.
- **The wider record set:** The controller database also holds everything that lives on the controller side instead of per model: {ref}`users <user>`, their {ref}`access levels <user-access-levels>`, {ref}`clouds <cloud>` and {ref}`credentials <credential>`, {ref}`SSH keys <ssh-key>`, {ref}`secret backends <secret-backend>`, leases, and the {ref}`migration <the-model-migration>` bookkeeping. The per-model {ref}`databases <database>` hold the rest.

(the-controller-watchers)=
## Controller change streams

Juju does not poll for changes to the controller. Its services publish change streams over the controller records, and the controller's own workers and the agents that need to know subscribe to them.

- **Controller configuration:** The controller config record and the controller record; the controller's own workers reconcile on it.
- **Controller nodes:** The HA cluster's membership.
- **The API addresses:** The addresses the controller's nodes serve on, so agents can re-orient as the cluster changes.

Every watcher fires once immediately when it is created, the initial query being the baseline snapshot, and again on each qualifying change. See {ref}`the watcher pattern <watchers>`.

(the-controllers-machinery)=
## Controller operations

Once running, every verb on a controller is a call over the controller API, and the access level decides who may make it:

- **Login access:** Reading controller state, for example listing controllers.
- **Superuser access:** Changing the controller's configuration or destroying it.

(the-controller-operations)=
### Controller bootstrap

Bootstrap turns an empty cloud into a running control plane, in two halves: the client's, then the controller's (for example, with `juju bootstrap`).

1. **The client's half.** The client authenticates directly against the target cloud; there is no controller to gate access yet. It provisions a host machine, installs `jujuagentd` and waits. It records the new controller in its own store.
2. **The controller's half.** The host machine starts its controller agent, which runs the API server and the database in-process, then reports the API ready. The controller declares itself: it names itself from the `controller-name` key of the {ref}`controller configuration <list-of-controller-configuration-keys>`, and it declares its own {ref}`application <application>`, named `controller`, in the controller model it creates.

The result is one controller, one model and no applications besides the controller's own: the controller machine running `jujuagentd`, with the API server and Dqlite in-process. Bootstrap is also the controller record's creator: the initialised database gets the controller record, the `admin` {ref}`user <user>` and its superuser access.

```{ibnote}
See more: {ref}`bootstrap-a-controller`, {ref}`manage-controllers`
```

### Controller configuration

The controller configuration is the controller's key/value record, read and written through the controller config service and watched for changes by the controller's own workers, which subscribe to the {ref}`controller-configuration change stream <the-controller-watchers>`. Changing it needs superuser access. See {ref}`the list of controller configuration keys <list-of-controller-configuration-keys>`.

### Controller high availability

High availability means running more than one controller node. The nodes are the HA cluster's membership; adding a node brings another controller agent into the controller's Dqlite cluster, and the controller-nodes and API-addresses {ref}`change streams <the-controller-watchers>` report the change. Juju 4.1 has no `juju enable-ha` command; the `EnableHA` call on the controller API adds and removes the nodes.

### Controller removal

Destroying a controller needs superuser access. The cloud the controller was bootstrapped on cannot be removed while the controller stands (see {ref}`cloud <cloud>`).

