[private]
default:
  @just --list container

start:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.stop=true \
    --filter status=exited \
      | xargs -r docker restart > /dev/null

stop:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.stop=true \
      | xargs -r docker stop > /dev/null

exec:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.exec \
      | xargs -r -I {} \
        docker inspect --format '{{{{.Id}}}} {{{{index .Config.Labels "docker-restic.container.exec"}}}}' {} \
      | xargs -r -n2 \
        /bin/sh -c 'docker exec $0 /bin/sh -c "$1"'
