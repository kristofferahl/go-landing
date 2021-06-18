# go-landing

A minimalistic, templated landing page web server written in go.

![go-landing screenshot](screenshot.png)

[![Docker Automated build](https://img.shields.io/docker/cloud/automated/kristofferahl/go-landing.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/go-landing/)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/kristofferahl/go-landing.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/go-landing/)
[![MicroBadger Size](https://img.shields.io/microbadger/image-size/kristofferahl/go-landing.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/go-landing/)
[![Docker Pulls](https://img.shields.io/docker/pulls/kristofferahl/go-landing.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/go-landing/)


## Build
```bash
go build .
```

## Run
```bash
go run .
```

## Configuration

| Environment variable | Description                                   | Default                   |
|----------------------|-----------------------------------------------|---------------------------|
| LANDING_TEMPLATE     | The path to the html template                 | templates/index.html.tmpl |
| LANDING_TITLE        | The title displayed on the landing page       | go-landing                |
| LANDING_DESCRIPTION  | The description displayed on the landing page | powered by //go:embed     |

```bash
export LANDING_TEMPLATE='templates/index.html.tmpl'
export LANDING_TITLE='go-landing'
export LANDING_DESCRIPTION='powered by //go:embed'
```

## Running in docker
```bash
docker build -t kristofferahl/go-landing:v1.0.0 .
docker run --rm -p 9000:9000 \
  -e LANDING_TITLE='go-landing' \
  kristofferahl/go-landing:v1.0.0
```
