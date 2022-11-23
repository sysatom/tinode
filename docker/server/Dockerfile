FROM alpine:latest

ARG VERSION=0.1
ENV VERSION=$VERSION

LABEL version=$VERSION

ENV TINODE_URL=""

RUN apk add --no-cache bash

WORKDIR /opt/app

# Fetch build from Github.
COPY app/static/ ./static
COPY app/templ/ ./templ
COPY app/tinode .
COPY docker/server/config.template ./tinode.conf
COPY docker/server/entrypoint.sh .

RUN chmod +x tinode && chmod +x entrypoint.sh
RUN mkdir subscribe
RUN ls -all
RUN pwd
RUN stat tinode
RUN stat entrypoint.sh

#ENTRYPOINT ./entrypoint.sh
CMD ["/opt/app/tinode"]

EXPOSE 6060
