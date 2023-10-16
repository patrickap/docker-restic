#!/bin/sh

(
  script_dir=$(dirname "$0")

  log -i "Attempting to acquire lock for task BACKUP with timeout of ${RESTIC_SCRIPT_LOCK_TIMEOUT} seconds ..."

  flock --timeout ${RESTIC_SCRIPT_LOCK_TIMEOUT} --exclusive 200

  log -i "Starting task BACKUP ..."

  ${script_dir}/container/stop.sh && \
  ${script_dir}/backup/execute.sh && \
  ${script_dir}/container/start.sh && \
  ${script_dir}/backup/prune.sh && \
  ${script_dir}/backup/check.sh

  if [ $? -ne 0 ]; then
    ${script_dir}/container/start.sh
    log -e "Completed task BACKUP with errors. Check log output above."
    exit 1
  else
    log -i "Completed task BACKUP. Check log output above."
  fi
) 200> "${RESTIC_ROOT}/restic.lock"
