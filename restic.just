just_flags := "--justfile ${DOCKER_RESTIC_DIR}/config/restic.just --working-directory ${DOCKER_RESTIC_DIR}"
restic_flags := "--repo ${DOCKER_RESTIC_DIR}/data/repository --password-file /run/secrets/restic-password"
rclone_flags := "--password-command: 'cat /run/secrets/rclone-password'"

default:
  @just {{just_flags}} --list

backup:
  just {{just_flags}} container_stop
  -restic backup /source {{restic_flags}}
  just {{just_flags}} container_start

backup_prune:
  restic forget {{restic_flags}} \
    --keep-daily 7 \
    --keep-weekly 4 \
    --keep-monthly 12 \
    --keep-yearly 2 \
    --group-by paths \
    --prune

backup_check:
  restic check {{restic_flags}} \
    --read-data

backup_check_subset value:
  restic check {{restic_flags}} \
    --read-data-subset: {{value}}

sync remotes:
  set -euo pipefail \
    && for remote in {{remotes}}; do \
        rclone sync ${DOCKER_RESTIC_DIR}/data/repository ${remote} {{rclone_flags}} \
          --stats 15m \
          --fast-list \
          --progress; \
      done

sync_check remotes:
  set -euo pipefail \
    && for remote in {{remotes}}; do \
        rclone check ${DOCKER_RESTIC_DIR}/data/repository ${remote} {{rclone_flags}} \
          --stats 15m \
          --fast-list \
          --progress; \
      done

dump:
  restic dump latest / {{restic_flags}} \
    | gpg \
      --batch \
      --yes \
      --symmetric \
      --cipher-algo AES256 \
      --passphrase-file /run/secrets/restic-password \
      -o ${DOCKER_RESTIC_DIR}/data/export/backup_$(date +'%Y-%m-%d_%H.%M.%S').tar.gpg

dump_prune:
  set -euo pipefail \
    && ls -t ${DOCKER_RESTIC_DIR}/data/export/backup_* \
      | tail +7 \
      | xargs -r rm -rf

dump_check:
  set -euo pipefail \
    && ls -t ${DOCKER_RESTIC_DIR}/data/export/backup_* \
      | head -1 \
      | xargs -r gpg \
        --batch \
        --yes \
        --decrypt \
        --passphrase-file /run/secrets/restic-password \
          | tar -tf - > /dev/null


[private]
container_start:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.stop=true \
    --filter status=exited \
      | xargs -r docker restart > /dev/null

[private]
container_stop:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.stop=true \
      | xargs -r docker stop > /dev/null

[private]
container_exec:
  set -euo pipefail \
  && docker ps \
    --quiet \
    --filter label=docker-restic.container.exec \
      | xargs -r -I {} \
        docker inspect --format '{{{{.Id}}}} {{{{index .Config.Labels "docker-restic.container.exec"}}}}' {} \
      | xargs -r -n2 \
        /bin/sh -c 'docker exec $0 /bin/sh -c "$1"'
