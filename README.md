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
- Docker Desktop

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

### Docker Desktop

As an alternative you could use <a href="https://docs.docker.com/desktop/kubernetes/">Kubernetes for Docker Desktop</a>.

## Istio

// TODO: add istio setup description and what deploy files we are using

### Installation

Download from Istio <a href="https://istio.io/latest/docs/setup/install/">installation guide</a>

To setup istio locally:
```shell
istioctl install --set profile=minimal
```

Note: before using Istio enable sidecar injection:
```shell
kubectl label namespace kubetools istio-injection=enabled
```

#### Accessing services outside the cluster

In order to access cluster services perform steps described in the <a href="https://istio.io/latest/docs/setup/getting-started/#determining-the-ingress-ip-and-ports">setup guide</a>

### Istio Dashboards

Istio service mesh has a variety of dashboards to monitor your cluster:
- Jaeger - for tracing
- Kiali - for service cluster architecture visualizations
- Prometheus and Grafana - for cluster metrics

#### Add Jaeger Dashboard

To enable Jaeger in Istio run:
```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/addons/jaeger.yaml
```

To enter Jaeger UI run:
```shell
istioctl dashboard jaeger
```

// TODO: add jaeger img

#### Add Kiali Dashboard

To enable Kiali in Istio run:
```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/addons/kiali.yaml
```

To enter Kiali run:
```shell
istioctl dashboard kiali
```

// TODO: add kiali img from istio-course

#### Add Prometheus and Grafana Dashboard

Grafana setup for Istio:
```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/addons/grafana.yaml
```

Prometheus setup for Istio:
```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/addons/prometheus.yaml
```

To enter Grafana run:
```shell
istioctl dashboard grafana
```

Visit `http://localhost:3000/dashboard/db/istio-mesh-dashboard` in your web browser

// TODO: add prom and grafana img from istio-course


## Deploying application on Local Kubernetes Cluster

Here described two option how to deploy Go application on kubernetes cluster locally:
- Kustomize
- Helm

### Kustomize

- `kustomize/bases` contains base configuration files for deployment, configmaps, services, etc.
- `kustomize/overlays` contains base overlay config for base files
  - you could add new overlay to this directory to kustomize values in configs

#### Deploying to kuber cluster via Kustomize

To get kustomize build:
```shell
kustomize build kustomize/overlays/local

# or

kubectl kustomize kustomize/overlays/local
```

To deploy a stask to the kubernetes cluster:
```shell
kubectl apply -k  kustomize/overlays/local
```

### Helm

// TODO: describe how to configure and deploy cluster with HELM
