#!/bin/sh

log -i "Pruning backup snapshots ..."
log -i "Keeping ${RESTIC_BACKUP_KEEP_DAILY} daily, ${RESTIC_BACKUP_KEEP_WEEKLY} weekly, ${RESTIC_BACKUP_KEEP_MONTHLY} monthly."
restic forget --keep-daily ${RESTIC_BACKUP_KEEP_DAILY} --keep-weekly ${RESTIC_BACKUP_KEEP_WEEKLY} --keep-monthly ${RESTIC_BACKUP_KEEP_MONTHLY} --prune

if [ $? -ne 0 ]; then
  log -w "Could not prune backup snapshots."
else
  log -s "Pruned backup snapshots."
fi
