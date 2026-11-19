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
    local exit_code=$?
    local cleanup_exit_code=0

    if [ "{{container_stop}}" = "true" ]; then
      log info "Restarting containers" "profile={{profile}}"
      just container start || cleanup_exit_code=$?

      if [ "$cleanup_exit_code" -ne 0 ]; then
        log error "Failed to restart containers" "profile={{profile}}" "exit_code=${cleanup_exit_code}"

        if [ "$exit_code" -eq 0 ]; then
          exit_code=$cleanup_exit_code
        fi
      fi
    fi

    exit "$exit_code"
  }

  trap cleanup EXIT

  if [ "{{container_exec}}" = "true" ]; then
    log info "Executing container commands" "profile={{profile}}"
    just container exec
  fi

  if [ "{{container_stop}}" = "true" ]; then
    log info "Stopping containers" "profile={{profile}}"
    just container stop
  fi

  log info "Executing Restic" "profile={{profile}}" "command=restic {{args}}"

  env_file="$(just profile get {{profile}} _meta.env-file)"
  if [ -f "${env_file}" ]; then
    log info "Loading Restic environment" "profile={{profile}}" "file=${env_file}"
    . "${env_file}"
  fi
  
  # Expand the profile arguments into positional parameters
  # This is done via Just to ensure proper shell quoting
  set -- {{ shell(f"just profile expand {{profile}} {{args}}") }}; \
  restic "$@"
