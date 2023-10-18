FROM restic/restic:0.16.0

ARG RESTIC_PASSWORD="" \
    RESTIC_ROOT="/srv/restic" \
    RESTIC_REMOTE="remote:restic" \
    RESTIC_BACKUP_KEEP_DAILY="7" \
    RESTIC_BACKUP_KEEP_WEEKLY="4" \
    RESTIC_BACKUP_KEEP_MONTHLY="12" \
    RESTIC_BACKUP_KEEP_YEARLY="2" \
    RESTIC_DUMP_KEEP_LAST="8" \
    RESTIC_LOCK_TIMEOUT="21600" \
    RESTIC_STOP_CONTAINER_LABEL="restic-stop=true"

ENV RESTIC_PASSWORD=$RESTIC_PASSWORD \
    RESTIC_ROOT=$RESTIC_ROOT \
    RESTIC_REMOTE=$RESTIC_REMOTE \
    RESTIC_BACKUP_KEEP_DAILY=$RESTIC_BACKUP_KEEP_DAILY \
    RESTIC_BACKUP_KEEP_WEEKLY=$RESTIC_BACKUP_KEEP_WEEKLY \
    RESTIC_BACKUP_KEEP_MONTHLY=$RESTIC_BACKUP_KEEP_MONTHLY \
    RESTIC_BACKUP_KEEP_YEARLY=$RESTIC_BACKUP_KEEP_YEARLY \
    RESTIC_DUMP_KEEP_LAST=$RESTIC_DUMP_KEEP_LAST \
    RESTIC_LOCK_TIMEOUT=$RESTIC_LOCK_TIMEOUT \
    RESTIC_STOP_CONTAINER_LABEL=$RESTIC_STOP_CONTAINER_LABEL \
    # internal only: do not change
    RESTIC_SOURCE="$RESTIC_ROOT/source" \
    RESTIC_TARGET="$RESTIC_ROOT/target" \
    RESTIC_REPOSITORY="$RESTIC_ROOT/target/repository" \
    RESTIC_EXPORT="$RESTIC_ROOT/target/export" \
    # add commands to PATH: do not change
    PATH="$RESTIC_ROOT/scripts:$PATH" \
    # set rclone config path: do not change
    RCLONE_CONFIG="$RESTIC_ROOT/config/rclone.conf"

COPY . $RESTIC_ROOT

RUN apk update \
    && apk add \
        docker-cli~=23.0.6 \
        rclone~=1.62.2 \
        flock~=2.38.1 \
        supercronic~=0.2.24 \
    && chmod -R 755 \
        $RESTIC_ROOT/scripts \
        $RESTIC_ROOT/entrypoint.sh \
    && mkdir -p \
        $RESTIC_SOURCE \
        $RESTIC_TARGET \
        $RESTIC_REPOSITORY \
        $RESTIC_EXPORT

ENTRYPOINT ["/bin/sh", "-c", "$RESTIC_ROOT/entrypoint.sh"]
