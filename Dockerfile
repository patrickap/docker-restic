FROM golang:1.23.0 as builder

WORKDIR /build
COPY . .

RUN go mod download \
    && go build -o ./bin/docker-restic

FROM restic/restic:0.17.0

ARG UID="1234" \
    GID="1234" \
    DOCKER_RESTIC_DIR="/srv/restic"

ENV UID=$UID \
    GID=$GID \
    DOCKER_RESTIC_DIR=$DOCKER_RESTIC_DIR \
    # set restic cache directory
    RESTIC_CACHE_DIR="$DOCKER_RESTIC_DIR/cache" \
    # set rclone config path
    RCLONE_CONFIG="$DOCKER_RESTIC_DIR/config/rclone.conf"

COPY --from=builder /build/bin /usr/bin
COPY --from=builder /build/entrypoint.sh /usr/bin/entrypoint.sh
COPY --from=builder /build/docker-restic.yml $DOCKER_RESTIC_DIR/config/docker-restic.yml
COPY --from=builder /build/docker-restic.cron $DOCKER_RESTIC_DIR/config/docker-restic.cron

RUN apk update \
    && apk add \
      docker-cli~=26.1.5 \
      rclone~=1.66.0 \
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

WORKDIR $DOCKER_RESTIC_DIR
ENTRYPOINT ["entrypoint.sh"]
CMD ["supercronic", "-passthrough-logs", "./config/docker-restic.cron"]
