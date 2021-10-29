#!/bin/bash

argocd app create kubetools --repo https://github.com/morzhanov/kuber-tools.git \
  --path kustomize/overlays/local \
  --sync-policy automatic \
  --dest-server http://your-kuber-cluster-url.svc
  --dest-namespace kubetools
