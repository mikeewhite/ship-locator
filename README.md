# 🚢 Ship Locator

![build status](https://github.com/mikeewhite/ship-locator/actions/workflows/backend.yml/badge.svg) [![codecov](https://codecov.io/gh/mikeewhite/ports-service/graph/badge.svg?token=BVGJXYFWCC)](https://codecov.io/gh/mikeewhite/ports-service)

![](https://github.com/mikeewhite/ship-locator/blob/main/images/demo.gif)

This repository represents a learning exercise in using the following technologies to build a ship location system that
exposes [AIS](https://en.wikipedia.org/wiki/Automatic_identification_system) data sourced from [aisstream.io](https://aisstream.io/) via a websocket API:

**Backend**
- Go (microservices)
- Kafka
- PostgreSQL
- Elasticsearch

**APIs**
 - gRPC
 - GraphQL

**Frontend**
- React
- Typescript
- Apollo Client
- Ant Design
- Google Maps API

**Observability**
- OpenTelemetry
- Jaeger
- Prometheus
- Grafana

## Usage

### Docker
Create a `.env` file in the root of this repo with the following content
```bash
# API key for websocket (get yours from https://aisstream.io/)
SHIPLOC_WEBSOCKETAPIKEY="YOUR-API-KEY" 

# API key for Google Maps Javascript API (see https://developers.google.com/maps/documentation/javascript/get-api-key)
REACT_APP_GOOGLE_MAPS_API_KEY="YOUR-API-KEY"
```
The applications can then be started via Docker using:
```bash
docker compose up -d
```

Once started the following services will be available:

| Name                | URL                   | Login                                                                     |
|---------------------|-----------------------|---------------------------------------------------------------------------|
| Dashboard           | http://localhost:3001 | -                                                                         |
| Grafana (metrics)   | http://localhost:3002 | `admin`/`admin`                                                           |
| Jaeger UI (tracing) | http://localhost:16686 | -                                                                         |
| pgAdmin (DB UI)     | http://localhost:5050 | `admin@admin.com`/`admin` (and `postgres` for saved server configuration) | 

### TODOs

- [ ] Move search query from GraphQL API to Ship Search Microservice and use a federated GraphQL server so to present the combined API as a single endpoint