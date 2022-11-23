FROM alpine:3.14

ARG VERSION=0.1
ENV VERSION=$VERSION

LABEL version=$VERSION

ENV TINODE_URL=""

WORKDIR /opt/tinode

RUN apk add --no-cache bash

# Fetch build from Github.
ADD tinode/app/static .
ADD tinode/app/templ .
ADD tinode/app/server .
ADD tinode/docker/server/config.template ./tinode.conf
COPY tinode/docker/server/entrypoint.sh .
RUN chmod +x server && chmod +x entrypoint.sh

ENTRYPOINT ./entrypoint.sh

EXPOSE 6060
