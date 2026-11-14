#!/bin/sh

set -- just \
  --justfile "${DOCKER_RESTIC_ETC_DIR}/docker-restic.conf" \
  --working-directory "${DOCKER_RESTIC_DATA_DIR}" \
  "${@}"

if [ "$(id -u)" -eq 0 ]; then
  exec su-exec restic "${@}"
else
  exec "${@}"
fi
