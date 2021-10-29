#!/bin/bash

export AMBASSADOR_LB_ENDPOINT=$(kubectl -n kubetools get svc ambassador -o "go-template={{range .status.loadBalancer.ingress}}{{or .ip .hostname}}{{end}}")
