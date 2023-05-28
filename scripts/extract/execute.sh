#!/bin/sh

log -i "Creating backup archive at '${RESTIC_ARCHIVE}' ..."
restic dump latest / > ${RESTIC_ARCHIVE}/archive_$(date +'%Y-%m-%d_%H.%M.%S').tar

if [ $? -ne 0 ]; then
  log -e "Could not create backup archive."
  exit 1
else
  log -s "Created backup archive."
fi
