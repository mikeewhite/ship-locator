FROM golang:1.20.3-alpine3.17 as builder

COPY . /build/
WORKDIR /build/

RUN go build -o collector ./cmd/collector
RUN go build -o service ./cmd/service
RUN go build -o search-service ./cmd/search-service

FROM alpine:3.17 as collector
WORKDIR /app
COPY --from=builder /build/collector ./

FROM alpine:3.17 as service
WORKDIR /app
COPY --from=builder /build/service ./

FROM alpine:3.17 as search-service
WORKDIR /app
COPY --from=builder /build/search-service ./
