#!/bin/sh

script_dir=$(dirname "$0")
error=0

log -i "Starting task CHECK ..."

${script_dir}/backup/check.sh || error=1
${script_dir}/dump/check.sh || error=1
${script_dir}/sync/check.sh || error=1

if [ $? -ne 0 ] || [ $error == 1 ]; then
  log -e "Completed task CHECK with errors. Check log output above."
  exit 1
else
  log -i "Completed task CHECK. Check log output above."
fi
