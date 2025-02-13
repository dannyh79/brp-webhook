# BRP Webhook Server

[![Tests](https://github.com/dannyh79/brp-webhook/actions/workflows/test.yml/badge.svg)](https://github.com/dannyh79/brp-webhook/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dannyh79/brp-webhook)](https://goreportcard.com/report/github.com/dannyh79/brp-webhook)

## Prereqs

- asdf
- make

## Getting Started

```sh
asdf install

make pre-flight

cp config.toml.example config.toml
# Then update the values in config.toml to your needs

# Start server; this exposes endpoints like POST localhost:8080/api/v1/callback
make
```

## Building

### Prereqs

- podman, or docker

```sh
podman build --file Dockerfile.linux.amd64 -t brp-webhook:linux-amd64 .
podman run --rm brp-webhook:linux-amd64 cat /app/bin/app > bin/app && chmod +x bin/app
```
