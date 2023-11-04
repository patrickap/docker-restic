#!/bin/sh

log -i "Creating backup snapshot of '${RESTIC_SOURCE_DIR}' ..."
restic -r ${RESTIC_REPOSITORY_DIR} backup ${RESTIC_SOURCE_DIR}

if [ $? -ne 0 ]; then
  log -e "Could not create backup snapshot."
  exit 1
else
  log -i "Created backup snapshot."
fi
