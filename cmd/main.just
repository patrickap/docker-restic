mod profile
mod container

set shell := ["/bin/sh", "-euo", "pipefail", "-uc"]

export JUST_JUSTFILE := justfile()

[private]
default:
  @just --list

[arg("profile", long="profile", short="p")]
[arg("container_stop", long="container-stop", short="s", value="true")]
[arg("container_exec", long="container-exec", short="x", value="true")]
exec profile container_exec="false" container_stop="false" *args:
  #!/usr/bin/env sh

  set -euo pipefail

  cleanup() {
    if [ "{{container_stop}}" = "true" ]; then
      just container start
    fi
  }

  trap cleanup EXIT

  if [ "{{container_exec}}" = "true" ]; then
    just container exec
  fi

  if [ "{{container_stop}}" = "true" ]; then
    just container stop
  fi

  log info "Executing Restic" "profile={{profile}}" "command=restic {{args}}"

  env_file="$(just profile get {{profile}} _meta.env-file)"
  if [ -f "${env_file}" ]; then
    . "${env_file}"
  fi

  restic $(just profile expand {{profile}} {{args}})
