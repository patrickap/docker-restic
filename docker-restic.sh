#!/bin/sh

if [ "$(id -u)" -eq 0 ]; then
  exec su-exec restic just --justfile ${DOCKER_RESTIC_DIR}/config/restic.just --working-directory ${DOCKER_RESTIC_DIR} "${@}"
else
  just --justfile ${DOCKER_RESTIC_DIR}/config/restic.just --working-directory ${DOCKER_RESTIC_DIR} "${@}"
fi
