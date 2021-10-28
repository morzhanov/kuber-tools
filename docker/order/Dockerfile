FROM golang:1.17-alpine as builder

WORKDIR /app

COPY ../../go.mod .
COPY ../../go.sum .

RUN echo $(pwd)
RUN go mod download

COPY ../../cmd/order/main.go ./cmd/order/main.go
COPY ../../api/ ./api/
COPY ../../internal/ ./internal

RUN go build -o ./bin/ ./cmd/order

# Final Stage - Stage 2
FROM alpine:3.11.5 as baseImage

WORKDIR /app

COPY --from=builder /app/bin/order /app/bin/order
COPY ../../config/.env.order ./config/.env.order

EXPOSE 3002

ENTRYPOINT [ "./bin/order" ]
