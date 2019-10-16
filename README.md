# airtrafficcontrol
[![Build Status](https://travis-ci.org/jsmithedin/airtrafficcontrol.svg?branch=master)](https://travis-ci.org/jsmithedin/airtrafficcontrol)

Slackbot to trigger deploys of Overmyhouse on merge to master

## Deployment
0. .env file containing slack bot key
```shell script
slackkey=BLAH
```
1. go build
2. Copy to /usr/local/airtrafficcontrol
3. Setup systemd
```shell script
cp airtrafficcontrol.service /etc/systemd/system
systemctl reload-daemon
systemctl enable airtrafficcontrol
systemctl start airtrafficcontrol
```
## Config
```shell script
url: URL of git repo
build: Build command
install: Install command
```
