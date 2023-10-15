#!/bin/sh

set -o pipefail

export="${RESTIC_ROOT}/target/export"

log -i "Pruning backup archives at '${export}' ..."
log -i "Keeping last ${RESTIC_DUMP_KEEP_LAST} backup archives."
ls -t ${export}/backup_* | tail +$((RESTIC_DUMP_KEEP_LAST+1)) | xargs -r rm

if [ $? -ne 0 ]; then
  log -w "Could not prune backup archives."
else
  log -i "Pruned backup archives."
fi
