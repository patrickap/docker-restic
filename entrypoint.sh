#!/bin/sh

echo "INF Starting container"

default_uid=$(id restic -u)
default_gid=$(id restic -g)

if [ ! "${UID}" = "${default_uid}" ] && [ -n "${UID}" ]; then
  echo "INF Changing UID from '${default_uid}' to '${UID}'"
  usermod -o -u "${UID}" restic
fi

if [ ! "${GID}" = "${default_gid}" ] && [ -n "${GID}" ]; then
  echo "INF Changing GID from '${default_gid}' to '${GID}'"
  groupmod -o -g "${GID}" restic
fi

if [ ! "${UID}" = "${default_uid}" ] || [ ! "${GID}" = "${default_gid}" ]; then
  echo "INF Changing ownership for '${DOCKER_RESTIC_DIR}' to '${UID}:${GID}'"
  chown -R restic:restic ${DOCKER_RESTIC_DIR}
fi

if capsh --has-b=cap_dac_read_search &> /dev/null; then
  echo "INF Setting necessary capabilities on binaries"
  setcap 'cap_dac_read_search=+ep' /usr/bin/restic
  setcap 'cap_dac_read_search=+ep' $DOCKER_RESTIC_DIR/bin/docker-restic
fi

echo "INF Running container as $(id restic)"
exec su-exec restic "${@}"