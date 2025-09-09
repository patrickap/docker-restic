#!/bin/sh

if [ "$(id -u)" -eq 0 ]; then
  exec su-exec restic just --justfile "${DOCKER_RESTIC_ETC_DIR}/docker-restic.conf" --working-directory "${DOCKER_RESTIC_DATA_DIR}" "${@}"
else
  just --justfile "${DOCKER_RESTIC_ETC_DIR}/docker-restic.conf" --working-directory "${DOCKER_RESTIC_DATA_DIR}" "${@}"
fi
