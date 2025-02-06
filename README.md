# BRP Webhook Server

## Prereqs

- asdf
- make

## Getting Started

```sh
asdf install

make pre-flight

# Start server; this exposes endpoints like POST localhost:8080/api/v1/callback
LINE_CHANNEL_SECRET="some-line-channel-secret" D1_GROUP_QUERY_ENDPOINT="http://example.com/api/v1/groups" make
```
