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
  @{{ if container_exec == "true" { "just container exec" } else { "" } }}
  @{{ if container_stop == "true" { "just container stop" } else { "" } }}

  @log info "Executing Restic" "profile={{profile}}" "command=restic {{args}}"

  -@env_file="$(just profile get {{profile}} _meta.env-file)"; \
    if [ -f "${env_file}" ]; then . "${env_file}"; fi; \
    restic $(just profile expand {{profile}} {{args}})

  @{{ if container_stop == "true" { "just container start" } else { "" } }}
