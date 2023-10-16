#!/bin/sh

(
  script_dir=$(dirname "$0")

  log -i "Attempting to acquire lock for task SYNC with timeout of ${RESTIC_SCRIPT_LOCK_TIMEOUT} seconds ..."

  flock --timeout ${RESTIC_SCRIPT_LOCK_TIMEOUT} --exclusive 200

  log -i "Starting task SYNC ..."

  ${script_dir}/sync/execute.sh && \
  ${script_dir}/sync/check.sh

  if [ $? -ne 0 ]; then
    log -e "Completed task SYNC with errors. Check log output above."
    exit 1
  else
    log -i "Completed task SYNC. Check log output above."
  fi
) 200> "${RESTIC_ROOT}/restic.lock"
