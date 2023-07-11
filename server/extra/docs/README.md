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
go run github.com/tinode/chat/server/extra/cmd/composer generator bot -name example -rule input,group,agent,command,condition,cron,form
go run github.com/tinode/chat/server/extra/cmd/composer generator vendor -name example

# Migrate cli
go run github.com/tinode/chat/server/extra/cmd/composer migrate import

# Migration file cli
go run github.com/tinode/chat/server/extra/cmd/composer migrate migration -name file_name
```
