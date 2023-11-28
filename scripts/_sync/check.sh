#!/bin/sh

log -i "Checking integrity of repository against remotes ..."

if [ -n "$RESTIC_RCLONE_REMOTES" ]; then
  IFS=','
  for remote in $RESTIC_RCLONE_REMOTES; do
    log -i "Checking against '${remote}' ..."
    rclone check ${RESTIC_ROOT_DIR}/backup/repository ${remote}

    if [ $? -ne 0 ]; then
      log -w "The remote data may be out of sync."
      log -e "Could not check integrity of '${remote}'."
      exit 1
    else
      log -i "The remote data seems fine."
      log -i "Checked integrity of '${remote}'."
    fi
  done
  unset IFS
else
  log -w "No remotes specified."
fi
