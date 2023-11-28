#!/bin/sh

log -i "Searching for containers labeled '${RESTIC_CONTAINER_EXEC_LABEL}' ..."
containers=$(docker ps -q --filter label="${RESTIC_CONTAINER_EXEC_LABEL}")

if [ -n "${containers}" ]; then
  log -i "Executing container commands ..."
  for container in ${containers}; do log -i "Executing '${container}' ..."; done
  command=$(docker inspect --format '{{index .Config.Labels "'${RESTIC_CONTAINER_EXEC_LABEL}'"}}' ${container})
  docker exec ${container} /bin/sh -c "${command}" > /dev/null

  if [ $? -ne 0 ]; then
    log -e "Could not execute container commands."
    exit 1
  else
    log -i "Executed container commands."
  fi
else
  log -w "No containers found. Possibly nothing to execute."
fi
