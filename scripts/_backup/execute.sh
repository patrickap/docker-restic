#!/bin/sh

log -i "Creating backup snapshot ..."
restic -r ${RESTIC_REPOSITORY_DIR} backup /source

if [ $? -ne 0 ]; then
  log -e "Could not create backup snapshot."
  exit 1
else
  log -i "Created backup snapshot."
fi
