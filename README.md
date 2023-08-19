# 🚢 Ship Locator

![build status](https://github.com/mikeewhite/ship-locator/actions/workflows/go.yml/badge.svg)

![](https://github.com/mikeewhite/ship-locator/demo.gif)

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
docker compose up
```