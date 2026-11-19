#!/bin/sh

log info "Starting container"

default_uid="$(id restic -u)"
default_gid="$(id restic -g)"

if [ ! "${UID}" = "${default_uid}" ] && [ -n "${UID}" ]; then
  log info "Changing UID" "from=${default_uid}" "to=${UID}"
  usermod -o -u "${UID}" restic
fi

if [ ! "${GID}" = "${default_gid}" ] && [ -n "${GID}" ]; then
  log info "Changing GID" "from=${default_gid}" "to=${GID}"
  groupmod -o -g "${GID}" restic
fi

if [ ! "${UID}" = "${default_uid}" ] || [ ! "${GID}" = "${default_gid}" ]; then
  log info "Changing ownership of directories" "from=${default_uid}:${default_gid}" "to=${UID}:${GID}"
  chown -R restic:restic \
    "${DOCKER_RESTIC_DATA_DIR}" \
    "${DOCKER_RESTIC_CONFIG_DIR}" \
    "${RESTIC_CACHE_DIR}"
fi

if capsh --has-b=cap_dac_read_search > /dev/null 2>&1; then
  log info "Applying capability to Restic" "capability=cap_dac_read_search"
  setcap 'cap_dac_read_search=+ep' /usr/bin/restic
  setcap 'cap_dac_read_search=+ep' /usr/bin/just
fi

for profile in $(docker-restic profile list); do
  if ! docker-restic exec --profile "${profile}" cat config > /dev/null 2>&1; then
    log info "Initializing Restic" "profile=${profile}"
    docker-restic exec --profile "${profile}" init
  fi
done

log info "Running container" $(id restic)
exec su-exec restic "$@"
