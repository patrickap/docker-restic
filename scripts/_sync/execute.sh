#!/bin/sh

log -i "Syncing repository to remotes ..."

if [ -n "$RESTIC_RCLONE_REMOTES" ]; then
  IFS=','
  for remote in $RESTIC_RCLONE_REMOTES; do
    log -i "Syncing to '${remote}' ..."
    rclone sync ${RESTIC_ROOT_DIR}/backup/repository ${remote} --progress --stats 15m

    if [ $? -ne 0 ]; then
      log -e "Could not sync repository to '${remote}'."
      exit 1
    else
      log -i "Synced repository to '${remote}'."
    fi
  done
  unset IFS
else
  log -w "No remotes specified."
fi
