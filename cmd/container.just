set shell := ["/bin/sh", "-euo", "pipefail", "-uc"]

stop_label := "docker-restic.container.stop=true"
exec_label := "docker-restic.container.exec"

[private]
default:
  @just --list container

start:
  @docker ps \
    --quiet \
    --filter label={{stop_label}} \
    --filter status=exited \
  | xargs -r docker restart > /dev/null

stop:
  @docker ps \
    --quiet \
    --filter label={{stop_label}} \
  | xargs -r docker stop > /dev/null

exec:
  @docker ps \
    --quiet \
    --filter label={{exec_label}} \
  | xargs -r -I {} \
      docker inspect --format '{{{{.Id}}}} {{{{index .Config.Labels "{{exec_label}}"}}}}' {} \
  | xargs -r -n2 \
      /bin/sh -c 'docker exec $0 /bin/sh -c "$1"'
