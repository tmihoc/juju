---
myst:
  html_meta:
    description: "SSH key reference: the public keys that users store on the controller, how a model authorises them, the operations that add, import, list and delete them, and the errors those raise."
---

(ssh-key)=
# SSH key

```{ibnote}
See also: {ref}`manage-ssh-keys`
```

```{ggarch}
:file: ../ssh-key.ggarch
:view: SSH key and its neighbours
:caption: A user owns many SSH keys, and each key can be authorised on many models. The controller stores the key once per user and records each authorisation separately.
:alt: A user owns many SSH keys. An SSH key is authorised on many models.
```

An **SSH key** is a public key in the [SSH](https://www.ssh.com/academy/ssh-keys) protocol that a {ref}`user <user>` stores on a {ref}`controller <controller>`. A model authorises a set of keys, and the people who hold the matching private keys can open SSH sessions to the model's machines through the controller (see {ref}`juju-ssh`). On the machines that Juju sets up with cloud-init, the account is named `ubuntu` and has passwordless sudo privileges, so authorising a key gives root access to those machines, and the decision deserves care.

The controller keeps keys in its own database, together with the users. A key has no life column. It comes into being when a user adds or imports it, and it goes away when the user deletes it from its last model or when the user is removed. Removing a model deletes only the authorisations of that model. A user can have many keys, and a key can be authorised on many models.

(ssh-key-data-model)=
## SSH key data model

```{ggarch}
:file: ../ssh-key.ggarch
:view: SSH key data model
:caption: The controller database stores a key once, in the record of its user, and links it to each model that authorises it. The fingerprint hash algorithm is a lookup. Juju stores every fingerprint as sha256.
:alt: A user public SSH key record belongs to one user and names one fingerprint hash algorithm. A model authorised keys record links one key to one model.
```

- **User public key (`user_public_ssh_key`):** The key's owner, its comment, its public key text, and its fingerprint with the algorithm that produced it.
  - *Rule:* A user stores the same key, or the same fingerprint, only once.
  - *Rule:* The comment of the controller's own key is reserved. At bootstrap on a machine cloud the controller generates an SSH key pair, keeps the private half, and stores the public half in the controller configuration key `system-ssh-keys`, so that it can reach the machines it manages. A user cannot add or import a key that carries the same comment.
  - *Related errors:*
    - **`public key already exists`**: *Trigger:* Adding a key that the user has already authorised on the model. *Remediation:* None needed; the key is in place. An import skips such keys without an error.
    - **`invalid public key`**: *Trigger:* Adding or importing a key that is not a valid public key. *Remediation:* Provide the public half of the key, in the format of an `authorized_keys` line.
    - **`key contains a reserved comment`**: *Trigger:* Adding or importing a key whose comment is the comment of the controller's own key. The client sees the message that the key contains a restricted comment. *Remediation:* Change the comment.
- **Fingerprint hash algorithm (`ssh_fingerprint_hash_algorithm`):** The allowed values are `md5` and `sha256`, and Juju computes and stores `sha256` for every key.
- **Model authorisation (`model_authorized_keys`):** Links one key to one model.
  - *Rule:* The set of keys that a model effectively authorises leaves out the keys of users who are removed or disabled. The controller derives this set in the view `v_model_authorized_keys`, and the keys stay stored.

(ssh-key-operations)=
## SSH key operations

A client lists, adds, imports and deletes keys through the controller API, authorised for the model it addresses. List, add and delete act on the authenticated user: the call accepts a user name, and the controller ignores it in favour of the identity that logged in. Import acts on the user named in the call.

- **List:** Returns the authorised keys of the authenticated user on the model. The caller needs read access to the model. A caller with admin access also qualifies. A client can ask for the full key text or only the fingerprint.
- **Add:** Stores one or more keys for the authenticated user and authorises them on the model.
  - *Authority:* Admin access to the model or superuser access to the controller, and the change block must not be on.
  - *Related errors:* `public key already exists`, `invalid public key` and `key contains a reserved comment`, each returned as its own error code (`public key already exists`, `invalid public key`, `invalid public key comment`).
- **Import:** Fetches the keys of an identity at an outside source and authorises them. The sources are Launchpad (`lp:<name>`) and GitHub (`gh:<name>`).
  - *Authority:* The same as add, and the change block must not be on.
  - *Rule:* An import skips keys that are already authorised on the model.
  - *Related errors:*
    - **`unknown import source`**: *Trigger:* A key id whose prefix is neither `lp:` nor `gh:`. *Remediation:* Use one of the two prefixes.
    - **`import subject not found`**: *Trigger:* The source has no such identity. *Remediation:* Check the spelling of the identity.
- **Delete:** Removes keys from the model, matched by fingerprint, public key text or comment. A user's key that no model authorises any more is deleted from the user's record too.
  - *Authority:* Admin access to the model or superuser access to the controller, and neither the remove block nor the change block may be on.
- **Removal of a user or a model:** Removing a user deletes the user's keys and their authorisations. Removing a model deletes the model's authorisations.

The command-line client offers `juju add-ssh-key`, `juju import-ssh-key`, `juju ssh-keys` (with `--full` for the key text) and `juju remove-ssh-key`.

(ssh-key-authentication)=
## SSH key authentication

The controller's SSH server checks a connection in two steps. It first authenticates the user: the key that the client presents must be one of the keys the controller stores for that user, whichever model it is authorised on. It then authorises the connection to the model that the client asked for: the user needs admin access to that model or superuser access to the controller (see {ref}`user access levels <user-access-levels>`), and the key that authenticated the user must be authorised on that model.

- *Related error:*
  - **`public key ... is not associated with model ...`**: *Trigger:* The user authenticates with a key that is stored but not authorised on the target model. *Remediation:* Authorise that key on the model, or present a different key.

Machine agents have a facade that serves the keys a model authorises, together with a stream of changes to them. The 4.1 code has no client of that facade.

## SSH key migration

A migration carries a model's authorised keys with the model. When a model arrives from an older controller, whose authorised keys travel in the model configuration, the import drops every key whose comment contains the comment of the controller's own key, so that the source controller keeps no SSH access to the migrated model.
