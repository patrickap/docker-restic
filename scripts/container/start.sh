#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_CONTAINER_STOP_LABEL}' and 'status=exited' ..."
containers=$(docker ps -q --filter label=${RESTIC_CONTAINER_STOP_LABEL} --filter "status=exited")

if [ -n "${containers}" ]; then
  log -i "Starting containers: ${containers} ..."
  docker restart ${containers} > /dev/null

  if [ $? -ne 0 ]; then
    log -w "Could not start containers."
  else
    log -s "Started containers."
  fi
else
  log -w "No containers found. Possibly already started."
fi
