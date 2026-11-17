FROM alpine:3.24

ARG UID="1234" \
    GID="1234"

ENV UID="${UID}" \
    GID="${GID}" \
    DOCKER_RESTIC_DATA_DIR="/var/lib/docker-restic" \
    DOCKER_RESTIC_CONFIG_DIR="/etc/docker-restic" \
    RESTIC_CACHE_DIR="/var/cache/restic"

COPY profiles.toml "${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"
COPY jobs.cron "${DOCKER_RESTIC_CONFIG_DIR}/jobs.cron"

COPY cmd /usr/local/lib/docker-restic/cmd
COPY docker-restic.sh /usr/local/bin/docker-restic
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN apk add --no-cache \
        docker-cli~=29.9.3 \
        restic~=0.18.1 \
        gnupg~=2.4.9 \
        shadow~=4.18.0 \
        libcap~=2.78 \
        su-exec~=0.3 \
        supercronic~=0.2.49 \
        gettext~=1.0 \
        yq~=4.53.3 \
    && apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community \
        just~=1.59.0 \
    && addgroup -S -g "${GID}" restic \
    && adduser -S -D -s /bin/sh -u "${UID}" -G restic restic \
    && mkdir -p \
        "${DOCKER_RESTIC_DATA_DIR}" \
        "${DOCKER_RESTIC_CONFIG_DIR}" \
        "${RESTIC_CACHE_DIR}" \
    && chown -R restic:restic \
        "${DOCKER_RESTIC_DATA_DIR}" \
        "${DOCKER_RESTIC_CONFIG_DIR}" \
        "${RESTIC_CACHE_DIR}" \
    && chmod +x \
        /usr/local/bin/docker-restic \
        /usr/local/bin/entrypoint.sh
WORKDIR "${DOCKER_RESTIC_DATA_DIR}"

ENTRYPOINT ["entrypoint.sh"]

CMD ["/bin/sh", "-c", "exec supercronic -passthrough-logs --no-reap ${DOCKER_RESTIC_CONFIG_DIR}/jobs.cron"]
