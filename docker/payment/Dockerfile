# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

WORKDIR /app

COPY ../../go.mod .
COPY ../../go.sum .

RUN echo $(pwd)
RUN go mod download

COPY ../../cmd/payment/main.go ./cmd/payment/main.go
COPY ../../config/.env.payment ./config/.env.payment
COPY ../../api/ ./api/
COPY ../../internal/ ./internal

RUN go build -o ./bin/ ./cmd/payment

EXPOSE 3003

ENTRYPOINT [ "./bin/payment" ]
