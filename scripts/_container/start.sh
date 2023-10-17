#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_STOP_CONTAINER_LABEL}' and 'status=exited' ..."
containers=$(docker ps -q --filter label=${RESTIC_STOP_CONTAINER_LABEL} --filter "status=exited")

if [ -n "$containers" ]; then
  log -i "Starting containers ..."
  for container in $containers
  do
    echo "'$container'"
  done
  docker restart ${containers} > /dev/null

  if [ $? -ne 0 ]; then
    log -w "Could not start containers."
  else
    log -i "Started containers."
  fi
else
  log -w "No containers found. Possibly already started."
fi
