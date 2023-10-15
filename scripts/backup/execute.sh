#!/bin/sh

source="${RESTIC_ROOT}/source"
repository="${RESTIC_ROOT}/target/repository"

log -i "Creating backup snapshot for '${repository}' of '${source}' ..."
restic -r ${repository} backup ${source}

if [ $? -ne 0 ]; then
  log -e "Could not create backup snapshot."
  exit 1
else
  log -i "Created backup snapshot."
fi
