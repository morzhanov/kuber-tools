#!/bin/bash

curl -sL https://raw.githubusercontent.com/crossplane/crossplane/release-1.5/install.sh | sh

kubectl crossplane install configuration registry.upbound.io/xp/getting-started-with-aws:v1.5.0

watch kubectl get pkg
