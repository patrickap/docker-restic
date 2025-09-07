FROM restic/restic:0.18.0

ARG UID="1234" \
    GID="1234" \
    DOCKER_RESTIC_DIR="/srv/restic"

ENV UID=$UID \
    GID=$GID \
    DOCKER_RESTIC_DIR=$DOCKER_RESTIC_DIR \
    DOCKER_RESTIC_BACKUP_KEEP_DAILY="7" \
    DOCKER_RESTIC_BACKUP_KEEP_WEEKLY="4" \
    DOCKER_RESTIC_BACKUP_KEEP_MONTHLY="12" \
    DOCKER_RESTIC_BACKUP_KEEP_YEARLY="2" \
    DOCKER_RESTIC_DUMP_KEEP_LAST="7" \
    # set restic cache directory
    RESTIC_CACHE_DIR="$DOCKER_RESTIC_DIR/cache" \
    # set rclone config path
    RCLONE_CONFIG="$DOCKER_RESTIC_DIR/config/rclone.conf"

COPY entrypoint.sh /usr/bin/entrypoint.sh
COPY docker-restic.sh /usr/local/bin/docker-restic
COPY restic.conf $DOCKER_RESTIC_DIR/config/restic.conf
COPY restic.cron $DOCKER_RESTIC_DIR/config/restic.cron

RUN apk update \
    && apk add \
      docker-cli~=27.3.1 \
      rclone~=1.68.2 \
      expect~=5.45.4 \
      gnupg~=2.4.7 \
      just~=1.37.0 \
      shadow~=4.16.0 \
      libcap~=2.71 \
      su-exec~=0.2 \
      supercronic~=0.2.33 \
      gettext~=0.22.5 \
    && addgroup -S -g $GID restic \
    && adduser -S -D -s /bin/sh -u $UID -G restic restic \
    && mkdir -p \
      $DOCKER_RESTIC_DIR \
      $DOCKER_RESTIC_DIR/data \
      $DOCKER_RESTIC_DIR/config \
      $DOCKER_RESTIC_DIR/cache \
      $DOCKER_RESTIC_DIR/tmp \
    && chown -R restic:restic $DOCKER_RESTIC_DIR \
    && chmod +x /usr/bin/entrypoint.sh  \
    && chmod +x /usr/local/bin/docker-restic


WORKDIR $DOCKER_RESTIC_DIR
ENTRYPOINT ["entrypoint.sh"]
CMD ["supercronic", "-passthrough-logs", "--no-reap", "./config/restic.cron"]
