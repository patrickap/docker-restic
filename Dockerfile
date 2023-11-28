FROM restic/restic:0.16.0

ARG UID="1000" \
    GID="1000" \
    RESTIC_PASSWORD \
    RESTIC_PASSWORD_FILE \
    RESTIC_ROOT_DIR="/srv/restic" \
    RESTIC_BACKUP_KEEP_DAILY="7" \
    RESTIC_BACKUP_KEEP_WEEKLY="4" \
    RESTIC_BACKUP_KEEP_MONTHLY="12" \
    RESTIC_BACKUP_KEEP_YEARLY="2" \
    RESTIC_DUMP_KEEP_LAST="8" \
    RESTIC_SYNC_REMOTE_MATCH="^restic-.*" \
    RESTIC_SYNC_REMOTE_DIR="restic" \
    RESTIC_LOCK_TIMEOUT="21600" \
    RESTIC_CONTAINER_STOP_LABEL="restic.container.stop" \
    RESTIC_CONTAINER_EXEC_LABEL="restic.container.exec" \
    RESTIC_CHOWN_ALL="false"

ENV UID=$UID \
    GID=$GID \
    RESTIC_PASSWORD=$RESTIC_PASSWORD \
    RESTIC_PASSWORD_FILE=$RESTIC_PASSWORD_FILE \
    RESTIC_ROOT_DIR=$RESTIC_ROOT_DIR \
    RESTIC_BACKUP_KEEP_DAILY=$RESTIC_BACKUP_KEEP_DAILY \
    RESTIC_BACKUP_KEEP_WEEKLY=$RESTIC_BACKUP_KEEP_WEEKLY \
    RESTIC_BACKUP_KEEP_MONTHLY=$RESTIC_BACKUP_KEEP_MONTHLY \
    RESTIC_BACKUP_KEEP_YEARLY=$RESTIC_BACKUP_KEEP_YEARLY \
    RESTIC_DUMP_KEEP_LAST=$RESTIC_DUMP_KEEP_LAST \
    RESTIC_SYNC_REMOTE_MATCH=$RESTIC_SYNC_REMOTE_MATCH \
    RESTIC_SYNC_REMOTE_DIR=$RESTIC_SYNC_REMOTE_DIR \
    RESTIC_LOCK_TIMEOUT=$RESTIC_LOCK_TIMEOUT \
    RESTIC_CONTAINER_STOP_LABEL=$RESTIC_CONTAINER_STOP_LABEL \
    RESTIC_CONTAINER_EXEC_LABEL=$RESTIC_CONTAINER_EXEC_LABEL \
    RESTIC_CHOWN_ALL=$RESTIC_CHOWN_ALL \
    # set internal variables
    RESTIC_REPOSITORY_DIR="$RESTIC_ROOT_DIR/backup/repository" \
    RESTIC_EXPORT_DIR="$RESTIC_ROOT_DIR/backup/export" \
    RESTIC_CONFIG_DIR="$RESTIC_ROOT_DIR/config" \
    RESTIC_SCRIPT_DIR="$RESTIC_ROOT_DIR/scripts" \
    # set restic cache directory
    RESTIC_CACHE_DIR="$RESTIC_ROOT_DIR/cache" \
    # set rclone config path
    RCLONE_CONFIG="$RESTIC_ROOT_DIR/config/rclone.conf" \
    # add commands to PATH
    PATH="$RESTIC_ROOT_DIR:$RESTIC_ROOT_DIR/scripts:$PATH"

COPY . $RESTIC_ROOT_DIR

RUN apk update \
    && apk add \
        docker-cli~=23.0.6 \
        rclone~=1.62.2 \
        flock~=2.38.1 \
        supercronic~=0.2.24 \
        shadow~=4.13 \
        su-exec~=0.2 \
    && mkdir -p \
        $RESTIC_ROOT_DIR \
        $RESTIC_REPOSITORY_DIR \
        $RESTIC_EXPORT_DIR \
        $RESTIC_CONFIG_DIR \
        $RESTIC_SCRIPT_DIR \
        $RESTIC_CACHE_DIR \
    && chmod -R 755 \
        $RESTIC_SCRIPT_DIR \
        $RESTIC_ROOT_DIR/entrypoint.sh \
        $RESTIC_ROOT_DIR/init.sh \
    && addgroup -S -g $GID restic \
    && adduser -S -H -D -s /bin/sh -u $UID -G restic restic \
    && chown -R restic:restic $RESTIC_ROOT_DIR

ENTRYPOINT ["entrypoint.sh"]
CMD ["init.sh"]
