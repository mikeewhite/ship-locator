FROM golang:1.20.3-alpine3.17 as builder

COPY . /build/
WORKDIR /build/

RUN go build -o collector ./cmd/collector

FROM alpine:3.17 as collector
WORKDIR /app
COPY --from=builder /build/collector ./
