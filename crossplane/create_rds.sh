#!/bin/bash

# NOTE: change secrets
AWS_PROFILE=default && echo -e "[default]\naws_access_key_id = $(aws configure get aws_access_key_id --profile $AWS_PROFILE)\naws_secret_access_key = $(aws configure get aws_secret_access_key --profile $AWS_PROFILE)" > creds.conf

# NOTE: provide creds files
kubectl create secret generic aws-creds -n crossplane-system --from-file=creds=./creds.conf

kubectl apply -f ./rds.yaml

kubectl get rdsinstance rdspostgres

kubectl describe secret aws-rdspostgres-conn -n kubetools
