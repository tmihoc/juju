---
myst:
  html_meta:
    description: "Set up Canonical Kubernetes cloud with Juju, including required services like DNS, ingress, local storage, and bootstrap configuration."
---

(cloud-canonical-k8s)=
# The Canonical Kubernetes cloud and Juju

This document describes details specific to using a Canonical Kubernetes cloud with Juju.

```{ibnote}
See more: [Canonical Kubernetes documentation](https://documentation.ubuntu.com/canonical-kubernetes/)
```

When using this cloud with Juju, it is important to keep in mind that it is a (1) Kubernetes cloud and (2) not some other cloud.

```{ibnote}
See more: {ref}`cloud-differences`
```

As the differences related to (1) are already documented generically in the rest of the docs, here we record just those that follow from (2).

## Requirements

### Services that must enabled

- `dns`
- `ingress` (technically not required, but you need it if you want to do anything meaningful)
- `local-storage`
- `network`

## Notes on `juju add-k8s`

Because Canonical Kubernetes uses its own `k8s kubectl` binary rather than the standard `kubectl`, you must pipe the kubeconfig explicitly:

```text
sudo k8s kubectl config view --raw | juju add-k8s <cloud name> --client
```

When piping via stdin, Juju cannot prompt interactively to ask whether to register the cloud on the client or a controller, so you must specify `--client` (or `--controller <name>`) explicitly.

Before you bootstrap:

- Create a custom `containerd` path, e.g., `export containerdBaseDir="/run/containerd-k8s"`.

- Resize `/run`, e.g., `sudo mount -o remount,size=10G /run`.

```{ibnote}
See more: https://github.com/canonical/k8s-snap/issues/1612
```

## Cloud-specific storage providers

```{ibnote}
See first: {ref}`storage-provider`
```

As for all Kubernetes clouds. See {ref}`storage-provider-kubernetes`.