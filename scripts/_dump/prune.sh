#!/bin/sh

set -o pipefail

log -i "Searching for backup archives to prune at '${RESTIC_EXPORT}' ..."
log -i "Keeping last ${RESTIC_DUMP_KEEP_LAST} backup archives."
prune_backups=$(ls -t ${RESTIC_EXPORT}/backup_* | tail +$((RESTIC_DUMP_KEEP_LAST+1)) | xargs -r echo)

if [ -n "$prune_backups" ]; then
  log -i "Pruning backup archives: '$(echo $prune_backups | xargs basename | xargs -r echo)' ..."
  rm -rf ${prune_backups}
  
  if [ $? -ne 0 ]; then
    log -w "Could not prune backup archives."
  else
    log -i "Pruned backup archives."
  fi
else
  log -w "Nothing to prune. Prune policy may not be fulfilled."
fi
