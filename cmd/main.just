mod profile
mod container

set shell := ["/bin/sh", "-euo", "pipefail", "-uc"]

export JUST_JUSTFILE := justfile()

[private]
default:
  @just --list

[arg("profile", long, short)]
[arg("suspend", long, short, value="true")]
exec profile suspend="false" *args:
  @{{ if suspend == "true" { "just container stop" } else { "" } }}
  @restic $(just profile expand {{profile}} {{args}})
  @{{ if suspend == "true" { "just container start" } else { "" } }}
