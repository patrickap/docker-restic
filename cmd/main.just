mod profile
mod container

export JUST_JUSTFILE := justfile()

[private]
default:
  @just --list

[arg("profile", long, short)]
[arg("suspend", long, short, value="true")]
exec profile suspend="false" *args:
  @{{ if suspend == "true" { "just container stop" } else { "" } }}

  @set -euo pipefail \
    && restic $(just profile expand {{profile}} {{args}})

  @{{ if suspend == "true" { "just container start" } else { "" } }}
