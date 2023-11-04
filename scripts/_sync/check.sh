#!/bin/sh

log -i "Checking integrity of repository '${RESTIC_REPOSITORY_DIR}' against remote '${RESTIC_RCLONE_REMOTE}' ..."
rclone check ${RESTIC_REPOSITORY_DIR} ${RESTIC_RCLONE_REMOTE}

if [ $? -ne 0 ]; then
  log -w "The remote data may be out of sync."
  log -e "Could not check integrity of remote data."
  exit 1
else
  log -i "The remote data seems fine."
  log -i "Checked integrity of remote data."
fi
