#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_CONTAINER_STOP_LABEL}' ..."
containers=$(docker ps -q --filter label=${RESTIC_CONTAINER_STOP_LABEL})

if [ -n "${containers}" ]; then
  log -i "Stopping containers: ${containers} ..."
  docker stop ${containers} > /dev/null

  if [ $? -ne 0 ]; then
    log -w "Could not stop containers."
  else
    log -i "Stopped containers."
  fi
else
  log -w "No containers found. Possibly already stopped."
fi
