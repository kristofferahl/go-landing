# go-landing

A minimalistic, templated landing page web server written in go.

## Build
```bash
go build .
```

## Run
```bash
go run .
```

## Customize
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
