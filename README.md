# 🚢 Ship Locator

![build status](https://github.com/mikeewhite/ship-locator/actions/workflows/go.yml/badge.svg)

## Usage

### Docker
Create a `.env` file in the root of this repo with the following content
```bash
# API key for websocket (get yours from https://aisstream.io/)
SHIPLOC_WEBSOCKETAPIKEY="YOUR-API-KEY" 
```
The applications can then be started via Docker using:
```bash
docker compose up
```