#!/bin/bash

docker build -t vladmorzhanov/kubetools_apigw:1.0.1 -f ./docker/apigw/Dockerfile .
docker push vladmorzhanov/kubetools_apigw:1.0.1

docker build -t vladmorzhanov/kubetools_order:1.0.5 -f ./docker/order/Dockerfile .
docker push vladmorzhanov/kubetools_order:1.0.5

docker build -t vladmorzhanov/kubetools_payment:1.0.5 -f ./docker/payment/Dockerfile .
docker push vladmorzhanov/kubetools_payment:1.0.5
