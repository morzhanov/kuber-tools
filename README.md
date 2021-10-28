# kuber-tools

Kubernetes popular tooling testing example.

TBD

## Local App Development

### Docker Compose

// TODO: describe how to setup dc

## Postman Collection

// TODO: describe and add screenshots about postman

## Local Cluster setup

Here described two option how to start kubernetes cluster locally:
- Minikube
- K3D

### Minikube

```shell
# install minikube
brew install minikube

# start cluster
minikube start

# check cluster
kubectl get nodes
```

Open service on minikube:

```shell
minikube service apigw -n kubetools

|-----------|-------|----------------|---------------------------|
| NAMESPACE | NAME  |  TARGET PORT   |            URL            |
|-----------|-------|----------------|---------------------------|
| kubetools | apigw | 3001-3001/3001 | http://192.168.49.2:30080 |
|-----------|-------|----------------|---------------------------|
🏃  Starting tunnel for service apigw.
|-----------|-------|-------------|------------------------|
| NAMESPACE | NAME  | TARGET PORT |          URL           |
|-----------|-------|-------------|------------------------|
| kubetools | apigw |             | http://127.0.0.1:51860 |
|-----------|-------|-------------|------------------------|
🎉  Opening service kubetools/apigw in default browser...
❗  Because you are using a Docker driver on darwin, the terminal needs to be open to run it.
```

### K3D

<a href="https://k3d.io/v5.0.1/">Installation</a>

```shell
# download k3d
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash

# create k3d cluster and open ports
k3d cluster create kubetools -p "30000-31652:30000-31652@server:0"

# check cluster
kubectl get nodes
```

## Kustomize

- `kustomize/bases` contains base configuration files for deployment, configmaps, services, etc.
- `kustomize/overlays` contains base overlay config for base files
  - you could add new overlay to this directory to kustomize values in configs

### Deploying to kuber cluster via Kustomize

To get kustomize build:
```shell
kustomize build kustomize/overlays/base

# or

kubectl kustomize kustomize/overlays/base
```

To deploy a stask to the kubernetes cluster:
```shell
kubectl apply -k  kustomize/overlays/base
```
