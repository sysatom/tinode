#!/bin/bash

# Make sure the system uses /etc/hosts when resolving domain names
# (needed for docker-compose's `extra_hosts` param to work correctly).
# See https://github.com/gliderlabs/docker-alpine/issues/367,
# https://github.com/golang/go/issues/35305 for details.
echo "hosts: files dns" > /etc/nsswitch.conf

# App
./server
