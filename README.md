# 🚢 Ship Locator

![build status](https://github.com/mikeewhite/ship-locator/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/mikeewhite/ports-service/graph/badge.svg?token=BVGJXYFWCC)](https://codecov.io/gh/mikeewhite/ports-service)

![](https://github.com/mikeewhite/ship-locator/blob/main/demo.gif)

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