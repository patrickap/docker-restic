#!/bin/sh
exec su-exec restic just --justfile ${DOCKER_RESTIC_DIR}/config/restic.just --working-directory ${DOCKER_RESTIC_DIR} "${@}"
