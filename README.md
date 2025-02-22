# BRP Webhook Server

[![Tests](https://github.com/dannyh79/brp-webhook/actions/workflows/test.yml/badge.svg)](https://github.com/dannyh79/brp-webhook/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dannyh79/brp-webhook)](https://goreportcard.com/report/github.com/dannyh79/brp-webhook)

## Prereqs

- asdf
- make

## Getting Started

```sh
asdf install

cp config.toml.example config.toml
# Then update the values in config.toml to your needs

# Build the binary then start server
# This exposes endpoints like POST localhost:8080/api/v1/callback
make
```

## Building for Linux AMD64 Platform

### Prereqs

- podman, or docker

```sh
make build-linux-amd64
```

## Developing

```sh
asdf install

cp config.toml.example config.toml
# Then update the values in config.toml to your needs

make pre-flight

# Available commands
make run
make fmt
make lint
make test
```
