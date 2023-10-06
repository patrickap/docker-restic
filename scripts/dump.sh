#!/bin/sh

script_dir=$(dirname "$0")

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
