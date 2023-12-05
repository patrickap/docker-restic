#!/bin/sh

log -i "Syncing repository to remotes ..."
remotes=$(rclone listremotes | grep "${RESTIC_SYNC_REMOTE_MATCH}" | tr -d :)
error=0

if [ -n "${remotes}" ]; then
  for remote in ${remotes}; do
    log -i "Syncing to '${remote}' ..."
    rclone sync ${RESTIC_REPOSITORY_DIR} ${remote}:${RESTIC_SYNC_REMOTE_DIR} --fast-list --progress --stats 15m

    if [ $? -ne 0 ]; then
      log -e "Could not sync to remote."
      error=1
    else
      log -i "Synced to remote."
    fi
  done
else
  log -w "Nothing synced. No remote found. Rclone may not be configured."
fi

if [ ${error} == 1 ]; then
  exit 1
fi
