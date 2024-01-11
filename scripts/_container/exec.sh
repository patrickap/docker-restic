#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_CONTAINER_EXEC_LABEL}' ..."
containers=$(docker ps -q --filter label="${RESTIC_CONTAINER_EXEC_LABEL}")
error=0

if [ -n "${containers}" ]; then
  log -i "Executing container commands ..."
  for container in ${containers}; do 
    command=$(docker inspect --format '{{index .Config.Labels "'${RESTIC_CONTAINER_EXEC_LABEL}'"}}' ${container})

    log -i "Executing '${command}' in '${container}' ..."
    docker exec ${container} /bin/sh -c "${command}" > /dev/null

    if [ $? -ne 0 ]; then
      log -e "Could not execute command."
      error=1
    else
      log -i "Executed command."
    fi
  done
else
  log -w "No containers found. Possibly nothing to execute."
fi

if [ ${error} == 1 ]; then
  exit 1
fi
