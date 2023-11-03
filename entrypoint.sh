#!/bin/sh

log -i "Initializing container ..."

if [ ! "$UID" = "$(id restic -u)" ] || [ ! "$GID" = "$(id restic -g)" ]; then
  if [ -n "$UID" ]; then
    log -i "Changing UID from $(id restic -u) to $UID."
    usermod -o -u "$UID" restic
  fi

  if [ -n "$GID" ]; then
    log -i "Updating GID from $(id restic -g) to $GID."
    groupmod -o -g "$GID" restic
  fi

  # change owner of restic root but exclude source directory
  find $RESTIC_ROOT -mindepth 1 -maxdepth 1 ! -name source -exec chown -R restic:restic {} \;
fi

if restic cat config &> /dev/null; then
  log -i "Skipping restic initialization. Repository already exists."
else
  restic -r ${RESTIC_REPOSITORY} init 2>&1

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

log -i "Running container as uid=$UID(restic) gid=$GID(restic)"
exec su-exec restic "$@"
