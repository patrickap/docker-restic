#!/bin/sh

log -i "Checking integrity of repository ..."
restic -r ${RESTIC_ROOT_DIR}/backup/repository check --read-data

if [ $? -ne 0 ]; then
  log -w "The repository may be corrupt."
  log -e "Could not check integrity of repository."
  exit 1
else
  log -i "The repository seems fine."
  log -i "Checked integrity of repository."
fi
