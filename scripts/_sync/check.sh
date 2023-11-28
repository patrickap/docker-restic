#!/bin/sh

log -i "Checking integrity of repository against remotes ..."
remotes=$(rclone listremotes | grep "${RESTIC_SYNC_REMOTE_MATCH}") | tr -d :
error=0

if [ -n "${remotes}" ]; then
  for remote in ${remotes}; do
    log -i "Checking against '${remote}' ..."
    rclone check ${RESTIC_REPOSITORY_DIR} ${remote}:${RESTIC_SYNC_REMOTE_DIR}

    if [ $? -ne 0 ]; then
      log -w "The remote data may be out of sync."
      log -e "Could not check integrity of '${remote}'."
      error=1
    else
      log -i "The remote data seems fine."
      log -i "Checked integrity of '${remote}'."
    fi
  done
else
  log -w "No matching remotes found. Rclone may not be configured."
fi

if [ ${error} == 1 ]; then
  exit 1
fi
