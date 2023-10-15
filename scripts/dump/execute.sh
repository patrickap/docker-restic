#!/bin/sh

repository="${RESTIC_ROOT}/target/repository"
export="${RESTIC_ROOT}/target/export"

log -i "Creating backup archive of '${repository}' at '${export}' ..."
backup=${export}/backup_$(date +'%Y-%m-%d_%H.%M.%S').tar
restic -r ${repository} dump latest / > ${backup}

if [ $? -ne 0 ]; then
  log -e "Could not create backup archive."
  exit 1
else
  log -i "Created backup archive '${backup}'."
fi
