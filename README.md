# BRP Webhook Server

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
