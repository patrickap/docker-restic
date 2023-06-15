#!/bin/sh

script_dir=$(dirname "$0")

log -i "Starting task EXTRACT ..."

${script_dir}/extract/execute.sh && \
${script_dir}/extract/prune.sh && \
${script_dir}/extract/check.sh

if [ $? -ne 0 ]; then
  log -e "Completed task EXTRACT with errors. Check log output above."
  exit 1
else
  log -i "Completed task EXTRACT. Check log output above."
fi
