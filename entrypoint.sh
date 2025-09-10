#!/bin/sh

echo "Starting container"

default_uid="$(id restic -u)"
default_gid="$(id restic -g)"

if [ ! "${UID}" = "${default_uid}" ] && [ -n "${UID}" ]; then
  echo "Changing UID from '${default_uid}' to '${UID}'"
  usermod -o -u "${UID}" restic
fi

if [ ! "${GID}" = "${default_gid}" ] && [ -n "${GID}" ]; then
  echo "Changing GID from '${default_gid}' to '${GID}'"
  groupmod -o -g "${GID}" restic
fi

if [ ! "${UID}" = "${default_uid}" ] || [ ! "${GID}" = "${default_gid}" ]; then
  echo "Changing ownership of directories to '${UID}:${GID}'"
  chown -R restic:restic "${DOCKER_RESTIC_HOME_DIR}" "${DOCKER_RESTIC_DATA_DIR}" "${DOCKER_RESTIC_CONFIG_DIR}" "${DOCKER_RESTIC_ETC_DIR}" "${DOCKER_RESTIC_CACHE_DIR}"
fi

if capsh --has-b=cap_dac_read_search &> /dev/null; then
  echo "Applying 'cap_dac_read_search' capability to Restic binaries"
  setcap 'cap_dac_read_search=+ep' /usr/bin/restic
  setcap 'cap_dac_read_search=+ep' /usr/bin/just
fi

if [ ! -d "${DOCKER_RESTIC_DATA_DIR}/repository" ] && [ -f /run/secrets/restic-password ]; then
  echo "Initializing Restic"
  su-exec restic restic init $(docker-restic --evaluate restic_flags | envsubst)
  # Create directory for exported archives (repository dumps).
  su-exec restic mkdir -p "${DOCKER_RESTIC_DATA_DIR}/export"
fi

if [ ! -f "${DOCKER_RESTIC_CONFIG_DIR}/rclone.conf" ] && [ -f /run/secrets/rclone-password ]; then
  echo "Initializing Rclone"
  # Rclone does not provide a non-interactive method to encrypt the configuration file via CLI. 
  # Therefore, the `expect` tool is used to automate the interactive encryption process.
  su-exec restic expect <<EOF
set timeout 1
spawn rclone config encryption set
expect "password:" { send -- "[exec cat /run/secrets/rclone-password]\r" }
expect "password:" { send -- "[exec cat /run/secrets/rclone-password]\r" }
expect eof
EOF
fi

# Change working directory
cd "${DOCKER_RESTIC_CONFIG_DIR}"

echo "Running container as $(id restic)"
exec su-exec restic "${@}"
