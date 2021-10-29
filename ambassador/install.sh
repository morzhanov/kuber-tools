#!/bin/bash

kubectl apply -f https://app.getambassador.io/yaml/edge-stack/latest/aes-crds.yaml && \
kubectl wait --for condition=established --timeout=90s crd -lproduct=aes && \
kubectl apply -f https://app.getambassador.io/yaml/edge-stack/latest/aes.yaml && \
kubectl -n kubetools wait --for condition=available --timeout=90s deploy -lproduct=aes
