FROM golang:1.17-alpine as builder

WORKDIR /app

COPY ../../go.mod .
COPY ../../go.sum .

RUN echo $(pwd)
RUN go mod download

COPY ../../cmd/payment/main.go ./cmd/payment/main.go
COPY ../../api/ ./api/
COPY ../../internal/ ./internal

RUN go build -o ./bin/ ./cmd/payment

# Final Stage - Stage 2
FROM alpine:3.11.5 as baseImage

WORKDIR /app

COPY --from=builder /app/bin/payment /app/bin/payment
COPY ../../config/.env.payment ./config/.env.payment
COPY ../../internal/payment/migrations ./internal/payment/migrations

EXPOSE 3003

ENTRYPOINT [ "./bin/payment" ]
