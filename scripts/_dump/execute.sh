#!/bin/sh

log -i "Creating backup archive ..."
backup=${RESTIC_ROOT_DIR}/backup/export/backup_$(date +'%Y-%m-%d_%H.%M.%S').tar
restic -r ${RESTIC_ROOT_DIR}/backup/repository dump latest / > ${backup}

if [ $? -ne 0 ]; then
  log -e "Could not create backup archive."
  exit 1
else
  log -i "Created backup archive '$(basename "$backup")'."
fi
