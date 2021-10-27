# kuber-tools

Kubernetes popular tooling testing example.

TBD

## K3D Setup

<a href="https://k3d.io/v5.0.1/">Installation</a>

```bash
# download k3d
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash

# create k3d cluster
k3d cluster create kubetools

# check cluster
kubectl get nodes
```

## Kustomize

