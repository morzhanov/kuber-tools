FROM golang:1.17-alpine as builder

WORKDIR /app

COPY ../../go.mod .
COPY ../../go.sum .

RUN echo $(pwd)
RUN go mod download

COPY ../../cmd/apigw/main.go ./cmd/apigw/main.go
COPY ../../api/ ./api/
COPY ../../internal/ ./internal

RUN go build -o ./bin/ ./cmd/apigw

# Final Stage - Stage 2
FROM alpine:3.11.5 as baseImage

WORKDIR /app

COPY --from=builder /app/bin/apigw /app/bin/apigw
COPY ../../config/.env.apigw ./config/.env.apigw

EXPOSE 3001

ENTRYPOINT [ "./bin/apigw" ]
