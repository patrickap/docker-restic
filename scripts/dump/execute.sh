#!/bin/sh

log -i "Creating backup archive at '${RESTIC_EXPORT}' ..."
restic dump latest / > ${RESTIC_EXPORT}/backup_$(date +'%Y-%m-%d_%H.%M.%S').tar

if [ $? -ne 0 ]; then
  log -e "Could not create backup archive."
  exit 1
else
  log -i "Created backup archive."
fi
