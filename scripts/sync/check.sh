#!/bin/sh

repository="${RESTIC_ROOT}/target/repository"

log -i "Checking integrity of repository '${repository}' against remote '${RESTIC_REMOTE}' ..."
rclone check ${repository} ${RESTIC_REMOTE}

if [ $? -ne 0 ]; then
  log -w "The remote data may be out of sync."
  log -e "Could not check integrity of remote data."
  exit 1
else
  log -i "The remote data seems fine."
  log -i "Checked integrity of remote data."
fi
