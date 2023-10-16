#!/bin/sh

log -i "Creating backup snapshot of '${RESTIC_SOURCE}' ..."
restic -r ${RESTIC_REPOSITORY} backup ${RESTIC_SOURCE}

if [ $? -ne 0 ]; then
  log -e "Could not create backup snapshot."
  exit 1
else
  log -i "Created backup snapshot."
fi
