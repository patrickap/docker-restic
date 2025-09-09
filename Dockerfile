FROM alpine:3.22

ARG UID="1234" \
    GID="1234" \
    DOCKER_RESTIC_HOME_DIR="/var/lib/docker-restic" \
    DOCKER_RESTIC_DATA_DIR="${DOCKER_RESTIC_HOME_DIR}/data" \
    DOCKER_RESTIC_CONFIG_DIR="${DOCKER_RESTIC_HOME_DIR}/config" \
    DOCKER_RESTIC_ETC_DIR="/etc/docker-restic" \
    DOCKER_RESTIC_CACHE_DIR="/var/cache/docker-restic"

ENV UID="${UID}" \
    GID="${GID}" \
    DOCKER_RESTIC_HOME_DIR="${DOCKER_RESTIC_HOME_DIR}" \
    DOCKER_RESTIC_DATA_DIR="${DOCKER_RESTIC_DATA_DIR}" \
    DOCKER_RESTIC_CONFIG_DIR="${DOCKER_RESTIC_CONFIG_DIR}" \
    DOCKER_RESTIC_ETC_DIR="${DOCKER_RESTIC_ETC_DIR}" \
    DOCKER_RESTIC_CACHE_DIR="${DOCKER_RESTIC_CACHE_DIR}" \
    DOCKER_RESTIC_BACKUP_KEEP_DAILY="7" \
    DOCKER_RESTIC_BACKUP_KEEP_WEEKLY="4" \
    DOCKER_RESTIC_BACKUP_KEEP_MONTHLY="12" \
    DOCKER_RESTIC_BACKUP_KEEP_YEARLY="2" \
    DOCKER_RESTIC_DUMP_KEEP_LAST="7" \
    # set restic cache directory
    RESTIC_CACHE_DIR="${DOCKER_RESTIC_CACHE_DIR}" \
    # set rclone config path
    RCLONE_CONFIG="${DOCKER_RESTIC_CONFIG_DIR}/rclone.conf"

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY docker-restic.sh /usr/local/bin/docker-restic
COPY docker-restic.conf "${DOCKER_RESTIC_ETC_DIR}/docker-restic.conf"
COPY docker-restic.cron "${DOCKER_RESTIC_CONFIG_DIR}/docker-restic.cron"

RUN apk add --no-cache \
      docker-cli~=28.3.3 \
      restic~=0.18.0 \
      rclone~=1.69.3 \
      expect~=5.45.4 \
      gnupg~=2.4.7 \
      just~=1.40.0 \
      shadow~=4.17.3 \
      libcap~=2.76 \
      su-exec~=0.2 \
      supercronic~=0.2.33 \
      gettext~=0.24.1 \
    && addgroup -S -g "${GID}" restic \
    && adduser -S -D -s /bin/sh -u "${UID}" -G restic restic \
    && mkdir -p \
      "${DOCKER_RESTIC_HOME_DIR}" \
      "${DOCKER_RESTIC_DATA_DIR}" \
      "${DOCKER_RESTIC_CONFIG_DIR}" \
      "${DOCKER_RESTIC_ETC_DIR}" \
      "${DOCKER_RESTIC_CACHE_DIR}" \
    && chown -R restic:restic \
      "${DOCKER_RESTIC_HOME_DIR}" \
      "${DOCKER_RESTIC_DATA_DIR}" \
      "${DOCKER_RESTIC_CONFIG_DIR}" \
      "${DOCKER_RESTIC_ETC_DIR}" \
      "${DOCKER_RESTIC_CACHE_DIR}" \
    && chmod +x /usr/local/bin/entrypoint.sh  \
    && chmod +x /usr/local/bin/docker-restic


WORKDIR "${DOCKER_RESTIC_HOME_DIR}"
ENTRYPOINT ["entrypoint.sh"]
CMD ["supercronic", "-passthrough-logs", "--no-reap", "./docker-restic.cron"]
