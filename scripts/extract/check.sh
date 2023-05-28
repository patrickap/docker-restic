#!/bin/sh

log -i "Searching for latest backup archive ..."
archive=$(ls -t ${RESTIC_ARCHIVE}/archive_* | head -1)

if [ -n "${archive}" ]; then
  log -i "Checking integrity of latest backup archive '${archive}' ..."
  echo ${archive} | xargs -r tar -tf > /dev/null

  if [ $? -ne 0 ]; then
    log -w "The backup archive may be corrupt."
    log -e "Could not check integrity of latest backup archive."
    exit 1
  else
    log -i "The backup archive seems fine."
    log -s "Checked integrity of latest backup archive."
  fi
else
  log -w "No archive found. Folder may be empty."
fi
