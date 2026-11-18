#!/bin/sh

set -- just --justfile /usr/local/lib/docker-restic/cmd/main.just "$@"

if [ "$(id -u)" -eq 0 ]; then
  exec su-exec restic "$@"
else
  exec "$@"
fi
