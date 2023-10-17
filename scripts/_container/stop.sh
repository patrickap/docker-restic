#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_STOP_CONTAINER_LABEL}' ..."
containers=$(docker ps -q --filter label=${RESTIC_STOP_CONTAINER_LABEL})

if [ -n "$containers" ]; then
  log -i "Stopping containers ..."
  for container in $containers
  do
    echo "'$container'"
  done
  docker stop ${containers} > /dev/null

  if [ $? -ne 0 ]; then
    log -w "Could not stop containers."
  else
    log -i "Stopped containers."
  fi
else
  log -w "No containers found. Possibly already stopped."
fi
