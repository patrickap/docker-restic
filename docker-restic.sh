#!/bin/sh

set -- just --justfile "${DOCKER_RESTIC_ETC_DIR}/cmd/main.just" "${@}"

if [ "$(id -u)" -eq 0 ]; then
  exec su-exec restic "${@}"
else
  exec "${@}"
fi
