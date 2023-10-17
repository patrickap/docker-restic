#!/bin/sh

set -o pipefail

log -i "Searching for backup archives to prune at '${RESTIC_EXPORT}' ..."
log -i "Keeping last ${RESTIC_DUMP_KEEP_LAST} backup archives."
backups=$(ls -t ${RESTIC_EXPORT}/backup_* | tail +$((RESTIC_DUMP_KEEP_LAST+1)) | xargs -r echo)

if [ -n "$backups" ]; then
  log -i "Pruning backup archives ..."
  for backup in $backups; do echo "'$(basename $backup)'"; done
  rm -rf ${backups}
  
  if [ $? -ne 0 ]; then
    log -w "Could not prune backup archives."
  else
    log -i "Pruned backup archives."
  fi
else
  log -w "Nothing to prune. Prune policy may not be fulfilled."
fi
