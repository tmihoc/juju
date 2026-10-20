---
myst:
  html_meta:
    description: "Juju relation reference: application connections through endpoints. These include regular relations, peer relations, subordinate relations, and cross-model relations, and involve relation databags."
---

(relation)=
# Relation

```{ibnote}
See also: {ref}`manage-relations`
```

In Juju, a **relation** is a connection an {ref}`application <application>` supports by virtue of having a particular {ref}`endpoint <application-endpoint>`.

## Relation taxonomy

```{ggarch}
:file: ../juju.ggarch
:view: Relation kinds
:caption: A relation joins two application endpoints, and the kind of relation depends on the applications. A regular relation joins two applications in one model. A peer relation joins the units of one application to one another. A subordinate relation joins a subordinate application to its principal application. A cross-model relation joins an application in one model to an application in another model.
:alt: A model holds Application A, Application B, Application D with units D/0 and D/1, and a principal and a subordinate application. Lines connect A to B (regular), D/0 to D/1 (peer), and the principal to the subordinate (subordinate). A line labelled cross-model joins Application B to Application C in a second, separate model.
```

(peer-relation)=
### Peer relation

A **peer** relation is a relation that an application has to itself (i.e., its units respond to one another) automatically by virtue of having a `peers` endpoint.

Because every relation results in the creation of unit and application databags in Juju's database, peer relations are sometimes used by charm authors as a way to persist charm data. When the application has multiple units, peer relations are also the mechanism behind {ref}`high availability <high-availability>`.

(non-peer-relation)=
### Non-peer relation

A **non-peer** relation is a relation from one application to another, where the applications support the same endpoint interface and have opposite `provides` / `requires` endpoint roles.

(subordinate-relation)=
#### Subordinate relation

A **subordinate** relation is a {ref}`non-peer <non-peer-relation>` relation where one application is principal and the other subordinate.

A subordinate charm is by definition a charm deployed on the same machine as the principal charm it is intended to accompany. When you deploy a subordinate charm, it appears in your Juju model as an application with no unit. The subordinate relation helps the subordinate application acquire a unit. The subordinate application then scales automatically when the principal application does, by virtue of this relation.

(the-implicit-juju-info-relation-endpoint)=
##### The implicit `juju-info` relation endpoint

Every application in a model implicitly provides an extra endpoint named `juju-info`, with role `provides`, interface `juju-info`, and global scope. The endpoint is supplied by Juju itself: it does not need to be (and cannot be) declared in the charm's `metadata.yaml` / `charmcraft.yaml`, and it does not appear in `juju info <charm>` or on the charm's Charmhub page. It is supported on every kind of charm, but is only useful for {ref}`machine charms <machine-charm>`, since it exists to allow {ref}`subordinate charms <subordinate-charm>` to attach to a principal that does not otherwise expose a suitable interface.

The `juju-info` endpoint is intended to be consumed by a {ref}`subordinate charm <subordinate-charm>` whose `metadata.yaml` / `charmcraft.yaml` declares an explicit `requires` endpoint with interface `juju-info` (typically with name `juju-info` and `scope: container`). The container scope is what causes the subordinate's unit to be co-located on the same machine as the principal's unit.

The `juju-info` endpoint is the standard mechanism used by general-purpose machine subordinates that need to ride along on every machine of a principal but do not have any application-specific data to exchange. For example, monitoring, logging, or system-administration agents such as [`ntp`](https://charmhub.io/ntp), [`telegraf`](https://charmhub.io/telegraf), [`nrpe`](https://charmhub.io/nrpe), and [`canonical-livepatch`](https://charmhub.io/canonical-livepatch).

The convention is for a subordinate to name its own `requires` endpoint `juju-info`, but any name can be used (for example: `info`, `general-info`); it just needs to use interface `juju-info`.

You integrate against the implicit endpoint the same way as any other endpoint:

```text
juju integrate <subordinate> <principal>
```

If the subordinate also has explicit endpoints whose interfaces match endpoints on the principal, those explicit endpoints take precedence over the implicit `juju-info` one when Juju resolves the relation. To force the implicit endpoint, name it explicitly on either side:

```text
juju integrate <subordinate>:<requires-endpoint> <principal>:juju-info
```

(non-subordinate-relation)=
#### Non-subordinate relation

A **non-subordinate** relation (aka 'regular') is a {ref}`non-peer <non-peer-relation>` relation where the applications are both principal.

(cross-model-relation)=
##### Cross-model relation

```{ibnote}
See also: {ref}`manage-relations`
```

A **cross-model** relation (aka 'CMR') is a {ref}`non-subordinate <non-subordinate-relation>` relation where the applications are on different models (+/- different controllers, +/- different clouds).

Cross-model relations enable, for example, scenarios where  your databases are hosted on bare metal, to take advantage of I/O performance, and your applications live within Kubernetes, to take advantage of scalability and application density.

If the network topology is anything other than flat, the Juju controllers will need to be bootstrapped with `--controller-external-ips`, `--controller-external-name`, or both, so that the controllers are able to communicate. Note that these config values can only be set at bootstrap time, and are read-only thereafter.

A cross-model relation has two sides: the offering side (aka "offerer") and the consume side (aka 'saas'). It does not make a difference which side of the relation (provider or requirer) is the offerer and which is the saas - the two are interchangeable. However, the endpoint type does influence on how juju sets up firewall rules: it is assumed that a requirer is the client and the provider is the server, so ports are opened on the provider side.

Note that application names are obfuscated (anonymised) to the offerer side:
- Applications that relate to the saas appear to the offerer as remote + token, e.g. `remote-76cd96ab50f146b284912afd1cc13a0e`.
- For the consumer, the remote app names is the saas name, e.g. `prometheus`.
##### Non-cross-model relation

A **non-cross-model** relation is a {ref}`non-subordinate <non-subordinate-relation>` relation where the applications are on  the same model.

## Relation identification

A relation is identified by a **relation ID** (assigned automatically by Juju; expressed in monotonically increasing numbers) or a **relation key** (derived from the endpoints, format: `application1:[endpoint] application2:[endpoint]`).

## Relation databag

When you create a relation between two applications, this results in the creation of relation databags. Databags are per relation and per application, and can be application-scoped or unit-scoped. Each unit involved in a relation gets a local copy of all the databags for that relation.

### Permissions around relation databags

```{ggarch}
:file: ../juju.ggarch
:slides: Non-peer relation databags | Peer relation databags
:slide-captions: In a non-peer relation (regular or subordinate), every unit reads and writes its own unit databag and reads the remote application's databags, and any unit reads the local application databag (`relation-get --app`) while only the leader unit (the one unit of an application that Juju elects to act for the whole application) writes it. The figure shows the leader unit. | In a peer relation, every unit reads and writes its own unit databag and reads every other unit's unit databag, and any unit reads the application databag while only the leader unit writes it. The permissions that a non-peer relation gives over the remote application turn inwards to the local application.
:alt: Two diagrams of databag access. In the first, a leader unit has thick arrows to the local application databag and its own unit databag, and dashed arrows to a remote application databag and a remote unit databag. In the second, a leader unit and another unit each have a thick arrow to their own unit databag and a dashed arrow to the other unit's databag, the leader unit has a thick arrow to the application databag, and the other unit has a dashed arrow to it.
```

While the relation is maintained,

- in a non-peer relation, whether regular or subordinate:
    - each unit can read and write to its own databag;
    - leader units can also read and write to the local application databag;
    - all units of an application can read all of the remote application's databags.
- in a peer relation:
    - each unit can read and write to its own databag;
    - leader units can also read and write to the application databag;
    - all units can read all of the application's databags. That is, whether leader or not, every unit can read its own unit databag as well as every other unit's unit databag as well as the application databag.

Note that, in peer relations, all permissions related to the remote application are turned inwards and become permissions related to the local application.

## Relation data model

```{ggarch}
:file: ../juju.ggarch
:view: Relation data model
:caption: The controller stores a relation in a `relation` record (with its life, scope, and a flag for suspension) and one `relation_endpoint` record for each application endpoint it joins. Each unit that has entered the relation has a `relation_unit` record that hangs on one endpoint. Databags are rows of key and value pairs: the unit databag of a unit is its `relation_unit_setting` rows, the application databag of an endpoint is its `relation_application_setting` rows, and a hash record for each bag (`relation_unit_settings_hash`, `relation_application_settings_hash`) lets watchers notice a change. Status, ingress and egress network addresses, the archive of settings from units that left, and the departure records hang on the relation itself.
:alt: An entity-relationship diagram of twelve tables. The relation table is on the left, with relation_status, the two network tables, relation_unit_departure, and relation_unit_setting_archive attached to it. The relation_endpoint table joins it to relation_application_setting and its hash table, and to relation_unit, which has relation_unit_setting and a hash table.
```
