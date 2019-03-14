# screen-monitor

## What's screen-monitor?
screen-monitor is a tool written in Golang that allows you to serve your server OS stats and screens stats
through API.

## Installation
* You need Golang Version 1.9.7.
```bash
$ git clone https://github.com/arsmine/screen-monitor
$ cd screen-monitor
$ make
```

## Usage
./screen-monitor --config example.json

Add screen names that you want to check to a json file in this format:
```json
{
	"activeScreen": ["screen-name"],
	"allowdIps": ["ip-address"]
}
```

## Command-Line Arguments

### required
* --config <config.json>

## Dependencies
|Package|
|:--|
|[go-osstat/memory](https://github.com/mackerelio/go-osstat/memory)|
|[go-osstat/cpu](https://github.com/mackerelio/go-osstat/cpu)|
|[go-osstat/uptime](https://github.com/mackerelio/go-osstat/uptime)|
|[go-osstat/disk](https://github.com/mackerelio/go-osstat/disk)|
|[go-osstat/network](https://github.com/mackerelio/go-osstat/network)|
|[go-osstat/loadavg](https://github.com/mackerelio/go-osstat/loadavg)|
|[gorilla/mux](https://github.com/gorilla/mux)|