# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

WORKDIR /app

COPY ../../go.mod .
COPY ../../go.sum .

RUN echo $(pwd)
RUN go mod download

COPY ../../cmd/apigw/main.go ./cmd/apigw/main.go
COPY ../../config/.env.order ./config/.env.order
COPY ../../api/ ./api/
COPY ../../internal/ ./internal

RUN go build -o ./bin/ ./cmd/apigw

EXPOSE 3001

ENTRYPOINT [ "./bin/apigw" ]
