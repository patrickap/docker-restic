#!/bin/sh

log -i "Searching for latest backup archive at '${RESTIC_EXPORT}' ..."
backup=$(ls -t ${RESTIC_EXPORT}/backup_* | head -1)

if [ -n "${backup}" ]; then
  log -i "Checking integrity of latest backup archive '${backup}' ..."
  echo ${backup} | xargs -r tar -tf > /dev/null

  if [ $? -ne 0 ]; then
    log -w "The backup archive may be corrupt."
    log -e "Could not check integrity of latest backup archive."
    exit 1
  else
    log -i "The backup archive seems fine."
    log -i "Checked integrity of latest backup archive."
  fi
else
  log -w "No archive found. Folder may be empty."
fi
