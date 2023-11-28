#!/bin/sh

if restic -r ${RESTIC_ROOT_DIR}/backup/repository cat config &> /dev/null; then
  log -i "Skipping restic initialization. Repository already exists."
else
  restic -r ${RESTIC_ROOT_DIR}/backup/repository init 2>&1

  if [ $? -ne 0 ]; then
    log -w "Could not initialize restic repository."
  else
    log -i "Initialized restic repository."
  fi
fi

rclone_remotes=$(rclone listremotes)

if [ -z "${rclone_remotes}" ]; then
  log -w "Rclone is not configured. Run 'rclone config'."
else
  log -i "Rclone is configured."
fi

restic_cron="${RESTIC_ROOT_DIR}/config/restic.cron"

log -i "Running container as $(id restic)."
supercronic -passthrough-logs "${restic_cron}"
