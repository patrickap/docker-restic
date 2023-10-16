#!/bin/sh

(
  script_dir=$(dirname "$0")

  log -i "Attempting to acquire lock for task DUMP with timeout of ${RESTIC_SCRIPT_LOCK_TIMEOUT} seconds ..."

  flock --timeout ${RESTIC_SCRIPT_LOCK_TIMEOUT} --exclusive 200

  log -i "Starting task DUMP ..."

  ${script_dir}/dump/execute.sh && \
  ${script_dir}/dump/prune.sh && \
  ${script_dir}/dump/check.sh

  if [ $? -ne 0 ]; then
    log -e "Completed task DUMP with errors. Check log output above."
    exit 1
  else
    log -i "Completed task DUMP. Check log output above."
  fi
) 200> "${RESTIC_ROOT}/restic.lock"
