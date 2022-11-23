FROM alpine:3.14

ARG VERSION=0.1
ENV VERSION=$VERSION

LABEL version=$VERSION

ENV TINODE_URL=""

WORKDIR /opt/tinode

RUN apk add --no-cache bash

# Fetch build from Github.
ADD app/static .
ADD app/templ .
ADD app/tinode .
ADD docker/server/config.template ./tinode.conf
RUN chmod +x tinode

CMD ["./tinode"]

EXPOSE 6060
