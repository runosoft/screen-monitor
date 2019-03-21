# screen-monitor

## What's screen-monitor?
screen-monitor is a tool written in Golang that allows you to serve your server OS stats' and screen stats'
through API.

## Installation
* You need Golang Version 1.9.7.
```bash
$ git clone https://github.com/arsmine/screen-monitor
$ cd screen-monitor
$ make
```

## Usage
`./screen-monitor --config example.json`

### config.json
* `interval`: Time interval(10s, 1m etc.) that how frequently you what to check screens.
* `listen`: IP address and port that you want to serve API.
* `activeScreen`: screen names that you want to check.
* `allowedIPs`: IP-addresses that can access your API.

```json
{
	"interval": "10s",
	"listen": "0.0.0.0:8080",
	"activeScreen": ["screen-name"],
	"allowedIPs": ["1.12.14.16"]
}
```

## Run as a Linux service
* add `screen-monitor.service` to `/etc/systemd/system/screen-monitor.service` (don't forget the fill the <...>)
* to send logs to `/var/log/screen-monitor.log`:
  - add `screen-monitor-log.conf` to `/etc/rsyslog.d/screen-monitor-log.conf` (don't forget the fill the <...>)
* `$ systemctl start screen-monitor.service`
* to start service on start-up:
  - `$ systemctl enable screen-monitor.service`


## Command-Line Arguments

### required
* `--config <config.json>`

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
