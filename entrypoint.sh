#!/bin/sh

echo "Starting container ..."

default_uid=$(id restic -u)
default_gid=$(id restic -g)

if [ ! "${UID}" = "${default_uid}" ] && [ -n "${UID}" ]; then
  echo "Changing UID from '${default_uid}' to '${UID}'."
  usermod -o -u "${UID}" restic
fi

if [ ! "${GID}" = "${default_gid}" ] && [ -n "${GID}" ]; then
  echo "Changing GID from '${default_gid}' to '${GID}'."
  groupmod -o -g "${GID}" restic
fi

if [ ! "${UID}" = "${default_uid}" ] || [ ! "${GID}" = "${default_gid}" ] || [ "${RESTIC_CHOWN_ALL}" = "true" ]; then
  echo "Changing ownership for '${DOCKER_RESTIC_DIR}' to '${UID}:${GID}'."
  chown -R restic:restic ${DOCKER_RESTIC_DIR}
fi

exec su-exec restic "${@}"
