# extra chatbot framework

## ENV

CHANNEL_PATH=/subscribe
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=123456
TINODE_URL=http://127.0.0.1:6060
DOWNLOAD_PATH=/download

## extra json config

> See extra.conf

## Dev tools

```shell

# Generator cli
go run github.com/tinode/chat/server/extra/cmd/generator -bot example -rule input,group,agent,command,condition,cron,form

# Migrate cli
go run github.com/tinode/chat/server/extra/cmd/migrate

# Migration file cli
go run github.com/tinode/chat/server/extra/cmd/migration example_name
```
