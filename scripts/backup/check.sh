#!/bin/sh

repository="${RESTIC_ROOT}/target/repository"

log -i "Checking integrity of repository '${repository}' ..."
restic -r ${repository} check --read-data

if [ $? -ne 0 ]; then
  log -w "The repository may be corrupt."
  log -e "Could not check integrity of repository."
  exit 1
else
  log -i "The repository seems fine."
  log -i "Checked integrity of repository."
fi
