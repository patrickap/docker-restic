#!/bin/sh

log -i "Initializing restic ..."
echo "$(restic version)"

if restic cat config &> /dev/null; then
  log -i "Skipping restic initialization. Repository already exists."
else
  restic init --repo=${RESTIC_REPOSITORY}

  if [ $? -ne 0 ]; then
    log -w "Could not initialize restic repository."
  else
    log -s "Initialized restic repository."
  fi
fi

# check rclone status
remote_name=$(echo ${RESTIC_REMOTE} | awk -F: '{print $1}')

if rclone listremotes | grep -q "$remote_name"; then
  log -s "The rclone remote '$remote_name' is configured."
else
  log -w "The rclone remote $remote_name is not configured. Run 'rclone config'."
fi

# configure cronjob
crontab /etc/restic/restic.cron

log -i "Running cron in foreground ..."
echo "$(crontab -l)"

crond -f
