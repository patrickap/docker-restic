#!/bin/sh

if restic cat config &> /dev/null; then
  log -i "Skipping restic initialization. Repository already exists."
else
  restic -r ${RESTIC_REPOSITORY_DIR} init 2>&1

  if [ $? -ne 0 ]; then
    log -w "Could not initialize restic repository."
  else
    log -i "Initialized restic repository."
  fi
fi

rclone_remote=$(echo ${RESTIC_RCLONE_REMOTE} | awk -F: '{print $1}')

if { rclone listremotes | grep -q "$rclone_remote"; } 2>&1; then
  log -i "The rclone remote '$rclone_remote' is configured."
else
  log -w "The rclone remote '$rclone_remote' is not configured. Run 'rclone config'."
fi

restic_cron="$RESTIC_ROOT_DIR/config/restic.cron"

log -i "Running container as $(id restic)."
supercronic -passthrough-logs "$restic_cron"
