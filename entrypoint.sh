#!/bin/sh

log -i "Starting container ..."

default_uid=$(id restic -u)
default_gid=$(id restic -g)

if [ ! "${UID}" = "${default_uid}" ] && [ -n "${UID}" ]; then
  log -i "Changing UID from '${default_uid}' to '${UID}'."
  usermod -o -u "${UID}" restic
fi

if [ ! "${GID}" = "${default_gid}" ] && [ -n "${GID}" ]; then
  log -i "Changing GID from '${default_gid}' to '${GID}'."
  groupmod -o -g "${GID}" restic
fi

if [ ! "${UID}" = "${default_uid}" ] || [ ! "${GID}" = "${default_gid}" ] || [ "${RESTIC_CHOWN_ALL}" = "true" ]; then
  log -i "Changing ownership for '${RESTIC_ROOT_DIR}' to '${UID}:${GID}'."
  chown -R restic:restic ${RESTIC_ROOT_DIR}
fi

exec su-exec restic "${@}"
