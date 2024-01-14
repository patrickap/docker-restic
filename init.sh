#!/bin/sh

echo "Running container as $(id restic)."
supercronic -passthrough-logs "${DOCKER_RESTIC_DIR}/default.cron"
