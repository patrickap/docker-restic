FROM alpine:3.20 as builder

WORKDIR /build
COPY . .

RUN apk update \
    && apk add \
      curl \
      unzip \
    && curl -L -o runr.zip https://github.com/patrickap/runr/archive/refs/tags/v0.1.0.zip \
    && unzip runr.zip -d . \
    && mkdir -p ./bin \
    && mv ./runr-0.1.0/build/runr ./bin

FROM restic/restic:0.17.0

ARG UID="1234" \
    GID="1234" \
    DOCKER_RESTIC_DIR="/srv/restic"

ENV UID=$UID \
    GID=$GID \
    DOCKER_RESTIC_DIR=$DOCKER_RESTIC_DIR \
    # set restic cache directory
    RESTIC_CACHE_DIR="$DOCKER_RESTIC_DIR/cache" \
    # set runr config path
    RUNR_CONFIG_DIR="$DOCKER_RESTIC_DIR/config" \
    # set rclone config path
    RCLONE_CONFIG="$DOCKER_RESTIC_DIR/config/rclone.conf"

COPY --from=builder /build/bin /usr/bin
COPY --from=builder /build/entrypoint.sh /usr/bin/entrypoint.sh
COPY --from=builder /build/restic.just $DOCKER_RESTIC_DIR/config/restic.just
COPY --from=builder /build/restic.cron $DOCKER_RESTIC_DIR/config/restic.cron
COPY --from=builder /build/rclone.conf $DOCKER_RESTIC_DIR/config/rclone.conf

RUN apk update \
    && apk add \
      docker-cli~=26.1.5 \
      rclone~=1.66.0 \
      expect~=5.45.4 \
      gnupg~=2.4.5 \
      just~=1.26.0 \
      shadow~=4.15.1 \
      libcap~=2.70 \
      su-exec~=0.2 \
      supercronic~=0.2.29 \
    && addgroup -S -g $GID restic \
    && adduser -S -H -D -s /bin/sh -u $UID -G restic restic \
    && mkdir -p \
      $DOCKER_RESTIC_DIR \
      $DOCKER_RESTIC_DIR/data \
      $DOCKER_RESTIC_DIR/config \
      $DOCKER_RESTIC_DIR/cache \
      $DOCKER_RESTIC_DIR/tmp \
    && chown -R restic:restic $DOCKER_RESTIC_DIR \
    && chmod +x /usr/bin/entrypoint.sh

# RUN echo "alias docker-restic='just --justfile $DOCKER_RESTIC_DIR/config/restic.just --working-directory $DOCKER_RESTIC_DIR'" >> /root/.bashrc

WORKDIR $DOCKER_RESTIC_DIR
ENTRYPOINT ["entrypoint.sh"]
CMD ["supercronic", "-passthrough-logs", "./config/restic.cron"]
