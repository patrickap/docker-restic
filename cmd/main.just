mod profile
mod container

export JUST_JUSTFILE := justfile()

[private]
default:
  @just --list

[env("RESTIC_CACHE_DIR", env("DOCKER_RESTIC_CACHE_DIR", ""))]
[arg("profile", long, short)]
[arg("suspend", long, short, value="true")]
exec profile suspend="false" *args:
  @{{ if suspend == "true" { "just container stop" } else { "" } }}
  restic $(just profile expand {{profile}} {{args}})
  @{{ if suspend == "true" { "just container start" } else { "" } }}
