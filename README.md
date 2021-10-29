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
<img src=""/>

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
<img src=""/>
Image taken from <a href="https://github.com/morzhanov/istio-2020">istio-2020</a> repo.

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
<img src=""/>
Image taken from <a href="https://github.com/morzhanov/istio-2020">istio-2020</a> repo.

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

## Rancher

In order to add Rancher to monitor the cluster we can deploy it as Docker image and connect our existing cluster.

Deploying Rancher:
```shell
docker run -d --restart=unless-stopped \
  -p 4000:80 -p 4001:443 \
  --privileged \
  rancher/rancher:latest
```

// TODO: write description and add images from desktop

## Crossplane

Crossplane goes beyond simply modelling infrastructure primitives as custom resources - it enables you to define new custom resources with schemas of your choosing.

<img src="https://www.percona.com/blog/wp-content/uploads/2021/05/crossplane-provider-sql.png" alt="crossplane"/>

As an example we can deploy AWS RDS Instance using crossplane, and it will be assigned to our cluster.

### Installing Crossplane

This should install Crossplane locally:
```shell
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/release-1.5/install.sh | sh
```

Install simple configuration for AWS (in prod it's better to use custom configuration):
```shell
kubectl crossplane install configuration registry.upbound.io/xp/getting-started-with-aws:v1.5.0
```

Wait until all packages become healthy:
```shell
watch kubectl get pkg
```

### Adding AWS RDS Instance to Kuber cluster

We can use Crossplane created Portgres instance instead of locally deployed by Kustomize file `kustomize/bases/postgresql`.

Using an AWS account with permissions to manage RDS databases:
```shell
AWS_PROFILE=default && echo -e "[default]\naws_access_key_id = $(aws configure get aws_access_key_id --profile $AWS_PROFILE)\naws_secret_access_key = $(aws configure get aws_secret_access_key --profile $AWS_PROFILE)" > creds.conf
```

Create a Provider Secret:
```shell
kubectl create secret generic aws-creds -n crossplane-system --from-file=creds=./creds.conf
```

The AWS provider supports provisioning an RDS instance via the RDSInstance managed resource it adds to Crossplane:
```shell
apiVersion: database.aws.crossplane.io/v1beta1
kind: RDSInstance
metadata:
  name: rdspostgres
spec:
  forProvider:
    region: us-east-1
    dbInstanceClass: db.t2.small
    masterUsername: masteruser
    allocatedStorage: 20
    engine: postgres
    engineVersion: "12"
    skipFinalSnapshotBeforeDeletion: true
  writeConnectionSecretToRef:
    namespace: kubetools
    name: aws-rdspostgres-conn
```

Note: RDSInstance is a Managed resource (MR). With Crossplane you could create a more complex Composite resources (XR):
<img src="https://crossplane.io/docs/v1.4/media/composition-xrs-and-mrs.svg" alt="rds"/>
You can review difference between managed and composite resources in the <a href="https://crossplane.io/docs/v1.4/concepts/composition.html">docs</a>

Creating the above instance will cause Crossplane to provision an RDS instance on AWS. You can view the progress with the following command:
```shell
kubectl get rdsinstance rdspostgres
```

When provisioning is complete, you should see `READY: True` in the output. You can take a look at its connection secret that is referenced under `spec.writeConnectionSecretToRef`:
```shell
kubectl describe secret aws-rdspostgres-conn -n kubetools
```
