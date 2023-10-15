#!/bin/sh

repository="${RESTIC_ROOT}/target/repository"

log -i "Syncing repository '${repository}' to remote '${RESTIC_REMOTE}' ..."
rclone sync ${repository} ${RESTIC_REMOTE} --progress --stats 15m

if [ $? -ne 0 ]; then
  log -e "Could not sync repository to remote."
  exit 1
else
  log -i "Synced repository to remote."
fi
