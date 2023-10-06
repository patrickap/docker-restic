#!/bin/sh

set -o pipefail

log -i "Pruning backup archives ..."
log -i "Keeping last ${RESTIC_DUMP_KEEP_LAST} backup archives."
ls -t ${RESTIC_ARCHIVE}/archive_* | tail +$((RESTIC_DUMP_KEEP_LAST+1)) | xargs -r rm

if [ $? -ne 0 ]; then
  log -w "Could not prune backup archives."
else
  log -i "Pruned backup archives."
fi
