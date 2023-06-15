#!/bin/sh

log -i "Syncing repository '${RESTIC_REPOSITORY}' to remote '${RESTIC_REMOTE}' ..."
rclone sync ${RESTIC_REPOSITORY} ${RESTIC_REMOTE} --progress --stats 15m

if [ $? -ne 0 ]; then
  log -e "Could not sync repository to remote."
  exit 1
else
  log -i "Synced repository to remote."
fi
