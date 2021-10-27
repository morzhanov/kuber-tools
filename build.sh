# !bash

docker build -t kubetools_apigw -f ./docker/apigw/Dockerfile .
docker build -t kubetools_order -f ./docker/order/Dockerfile .
docker build -t kubetools_payment -f ./docker/payment/Dockerfile .
