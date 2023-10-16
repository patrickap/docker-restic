#!/bin/sh

log -i "Initializing restic ..."
echo "$(restic version)"

if restic cat config &> /dev/null; then
  log -i "Skipping restic initialization. Repository already exists."
else
  restic -r "${RESTIC_ROOT}/target/repository" init 2>&1

  if [ $? -ne 0 ]; then
    log -w "Could not initialize restic repository."
  else
    log -i "Initialized restic repository."
  fi
fi

# check rclone status
remote_name=$(echo ${RESTIC_REMOTE} | awk -F: '{print $1}')

if { rclone listremotes | grep -q "$remote_name"; } 2>&1; then
  log -i "The rclone remote '$remote_name' is configured."
else
  log -w "The rclone remote '$remote_name' is not configured. Run 'rclone config'."
fi

log -i "Running cron in foreground ..."
cat "$RESTIC_ROOT/restic.cron"
supercronic -test "$RESTIC_ROOT/restic.cron"
supercronic -passthrough-logs "$RESTIC_ROOT/restic.cron"
