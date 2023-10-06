FROM restic/restic:0.16.0

ARG RESTIC_PASSWORD

ARG RESTIC_SOURCE="/source"
ARG RESTIC_REPOSITORY="/target/repository"
ARG RESTIC_ARCHIVE="/target/archive"
ARG RESTIC_REMOTE="remote:repository"

ARG RESTIC_BACKUP_KEEP_DAILY="7"
ARG RESTIC_BACKUP_KEEP_WEEKLY="4"
ARG RESTIC_BACKUP_KEEP_MONTHLY="12"
ARG RESTIC_BACKUP_KEEP_YEARLY="2"

ARG RESTIC_EXTRACT_KEEP_LAST="8"

ARG RESTIC_COMMAND_LOCK_FILE="/var/restic/cmd.lock"
ARG RESTIC_COMMAND_LOCK_TIMEOUT="21600"
ARG RESTIC_CONTAINER_STOP_LABEL="restic-stop=true"

ENV RESTIC_PASSWORD=$RESTIC_PASSWORD
ENV RESTIC_SOURCE=$RESTIC_SOURCE
ENV RESTIC_REPOSITORY=$RESTIC_REPOSITORY
ENV RESTIC_ARCHIVE=$RESTIC_ARCHIVE
ENV RESTIC_REMOTE=$RESTIC_REMOTE
ENV RESTIC_BACKUP_KEEP_DAILY=$RESTIC_BACKUP_KEEP_DAILY
ENV RESTIC_BACKUP_KEEP_WEEKLY=$RESTIC_BACKUP_KEEP_WEEKLY
ENV RESTIC_BACKUP_KEEP_MONTHLY=$RESTIC_BACKUP_KEEP_MONTHLY
ENV RESTIC_BACKUP_KEEP_YEARLY=$RESTIC_BACKUP_KEEP_YEARLY
ENV RESTIC_EXTRACT_KEEP_LAST=$RESTIC_EXTRACT_KEEP_LAST
ENV RESTIC_COMMAND_LOCK_FILE=$RESTIC_COMMAND_LOCK_FILE
ENV RESTIC_COMMAND_LOCK_TIMEOUT=$RESTIC_COMMAND_LOCK_TIMEOUT
ENV RESTIC_CONTAINER_STOP_LABEL=$RESTIC_CONTAINER_STOP_LABEL

# add commands to PATH for convenient execution
ENV PATH="/usr/local/bin/restic:$PATH"
# change rclone config path
ENV RCLONE_CONFIG="/etc/rclone/rclone.conf"

COPY ./restic.cron /etc/restic/restic.cron
COPY --chmod=0755 ./scripts/ /usr/local/sbin/restic/
COPY --chmod=0755 ./cmd/ /usr/local/bin/restic/
COPY --chmod=0755 ./entrypoint.sh /root/entrypoint.sh

RUN apk update && apk add \
    docker=23.0.6-r5 \
    rclone=1.62.2-r4 \
    flock=2.38.1-r8

RUN mkdir -p $RESTIC_REPOSITORY && \
    mkdir -p $RESTIC_ARCHIVE && \
    mkdir -p $(dirname $RESTIC_COMMAND_LOCK_FILE) && \
    mkdir -p $(dirname $RCLONE_CONFIG) && \
    touch $RESTIC_COMMAND_LOCK_FILE && \
    touch $RCLONE_CONFIG

ENTRYPOINT ["/root/entrypoint.sh"]