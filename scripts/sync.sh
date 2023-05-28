#!/bin/sh

script_dir=$(dirname "$0")

log -i "Starting task SYNC ..."

${script_dir}/sync/execute.sh && \
${script_dir}/sync/check.sh

if [ $? -ne 0 ]; then
  log -w "Completed task SYNC with errors. Check log output above."
  exit 1
else
  log -i "Completed task SYNC. Check log output above."
fi
