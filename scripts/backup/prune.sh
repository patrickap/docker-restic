#!/bin/sh

repository="${RESTIC_ROOT}/target/repository"

log -i "Pruning backup snapshots of '${repository}' ..."
log -i "Keeping ${RESTIC_BACKUP_KEEP_DAILY} daily, ${RESTIC_BACKUP_KEEP_WEEKLY} weekly, ${RESTIC_BACKUP_KEEP_MONTHLY} monthly, ${RESTIC_BACKUP_KEEP_YEARLY} yearly."
restic -r ${repository} forget --keep-daily ${RESTIC_BACKUP_KEEP_DAILY} --keep-weekly ${RESTIC_BACKUP_KEEP_WEEKLY} --keep-monthly ${RESTIC_BACKUP_KEEP_MONTHLY} --keep-yearly ${RESTIC_BACKUP_KEEP_YEARLY} --group-by paths --prune

if [ $? -ne 0 ]; then
  log -w "Could not prune backup snapshots."
else
  log -i "Pruned backup snapshots."
fi
