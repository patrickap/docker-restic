#!/bin/sh

log -i "Creating backup archive at '${RESTIC_EXPORT}' ..."
backup=${RESTIC_EXPORT}/backup_$(date +'%Y-%m-%d_%H.%M.%S').tar
restic -r ${RESTIC_REPOSITORY} dump latest / > ${backup}

if [ $? -ne 0 ]; then
  log -e "Could not create backup archive."
  exit 1
else
  log -i "Created backup archive '$(basename "$backup")'."
fi
