#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_CONTAINER_STOP_LABEL}=true' and 'status=exited' ..."
containers=$(docker ps -q --filter label="${RESTIC_CONTAINER_STOP_LABEL}=true" --filter "status=exited")

if [ -n "${containers}" ]; then
  log -i "Starting containers ..."
  for container in ${containers}; do log -i "Starting '${container}' ..."; done
  docker restart ${containers} > /dev/null

  if [ $? -ne 0 ]; then
    log -e "Could not start containers."
    exit 1
  else
    log -i "Started containers."
  fi
else
  log -w "No containers found. Possibly already started."
fi
