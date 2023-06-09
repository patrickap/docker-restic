# docker-restic

Docker-Restic is a Docker image that provides an easy way to use restic with additional features for container backups. It efficiently handles four key functionalities: backup snapshots, backup archives, remote synchronization and data integrity checks. The backup process is fully automated.

## Features

- **Easy Setup:** All data mounted at `/source` within Docker-Restic is backed up automatically to `/target`. This flexible setup allows you to define the specific directories and volumes you wish to include in your backups.
- **Backup Snapshots:** Docker-Restic performs daily snapshots using restic, allowing you to capture changes in your data efficiently.
- **Backup Archives:** Docker-Restic automatically exports a weekly tar archive file, providing a full backup of your data.
- **Remote Synchronization:** You have the option to enable a remote synchronization using rclone, which ensures that your backups are securely transferred to a remote location.
- **Integrity Checks**: Docker-Restic prioritizes the integrity of your backup data. It performs data integrity checks for all backup methods. These checks ensure that your backup data remains consistent and reliable, giving you peace of mind knowing that your valuable data is protected.
- **Fully Customizable:** Docker-Restic offers a high level of customization through various `ARG`s and `ENV`s that can be easily set or overwritten according to your requirements. These customization options provide the flexibility to adapt the backup process to your specific needs.
- **Various Extras:** Containers labeled with `restic-stop=true` are gracefully stopped before the backup process and restarted afterward, ensuring data consistency during the backup operation. To prevent concurrent access to backup resources, Docker-Restic utilizes a lockfile mechanism that effectively manages access and avoids conflicts.
- **Informative Logs**: Docker-Restic provides clear and easily comprehensible logs, making it effortless to monitor and troubleshoot the backup process. The logs are designed to present relevant information in a user-friendly format, enabling you to quickly identify any issues or track the progress of your backups.
- **Utility Commands**: Docker-Restic empowers you with the ability to perform manual backups and checks as needed. This feature allows you to take immediate backups of your container volumes or manually verify the integrity of existing backups.
- **All Restic Goodies**: Docker-Restic incorporates all the powerful features and capabilities of the restic backup tool. You can leverage restic's advanced functionalities, such as deduplication, encryption and data integrity checks to ensure robust and secure backups for your container volumes.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull the Docker-Restic image from the official Docker Hub repository:

```shell
docker pull patrickap/docker-restic:latest
```

2. Configure the necessary `ARG`s and `ENV`s to suit your backup requirements. Refer to the `Dockerfile` for a complete list of customization options.

3. Run the Docker-Restic container:

```bash
docker run -d --name docker-restic \
    -v /path/to/source:/source \
    -v /path/to/target:/target \
    -e RESTIC_PASSWORD=your-password \
    patrickap/docker-restic:latest
```

## Docker Compose

Here is a basic example how Docker-Restic can be used with docker compose.

```yml
version: '3.7'

services:
  restic:
    image: patrickap/docker-restic:latest
    environment:
      - RESTIC_PASSWORD=$BACKUP_PASSWORD
    init: true
    restart: always
    volumes:
      # backup destination
      - /path/to/backup:/target
      # volumes to backup
      - volume-1:/source/volume-1:ro
      - volume-2:/source/volume-2:ro
      # remote backup config
      - rclone-config:/etc/rclone
      # provide host information
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro

  volumes:
    volume-1:
    volume-2:
```

## Backup Volumes

By default mounted volumes inside `/source` are getting automatically backed up to `/target`. If you add custom volumes make sure to add them as read-only `:ro` for safety reasons. Also bind mount the backups to a custom location to be able to access them at any time.

```yml
docker-restic:
  volumes:
    - /path/to/backup:/target
    - volume-1:/source/volume-1:ro
    - volume-2:/source/volume-2:ro
    - volume-3:/source/volume-3:ro
```

## Remote Sync

Remote syncing of backups can be configured with `rclone`. Either bind mound the config into the container or run `rclone config` inside the `docker-restic` container.

```yml
docker-restic:
  volumes:
    - /path/to/rclone-config:/etc/rclone
```

**Note:**
If using Google Drive it is also recommended to create a [custom client-id](https://rclone.org/drive/#making-your-own-client-id) for better performance.

## Restore from Backup

To access or copy backups available on the remote host from another machine the command-line tool `scp` can be used.

```bash
scp username@<host_ip>:/path/to/source /path/to/target
```

To restore a backup it may be possible to use the official `restic restore` command with some additional setup. Otherwise a new Docker volume with the correct name must be created including the contents of the backup. After restarting the containers the data should be mounted and restored.

**Warning:**
If you're using Google Drive they may add back file extensions to encrypted files during the download or compression process which can result in a corrupted `restic` repository. To avoid this ensure to remove any added extensions inside the `repository/data` directory. An example of this would be a file at `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b` being downloaded as `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b.mp3`. In this example it is required to remove the automatically added extension `.mp3` to avoid repository corruption and be able to extract the backup.

```bash
# check restic repository
restic -r /path/to/repository check --read-data

# extract restic backup
restic -r /path/to/repository dump latest / > backup.tar

# untar the backup
tar -xvf backup.tar -C /tmp/backup

# stop the containers
docker stop <container_name>

# use a temporary once-off container to create the volume and copy the backup
docker run -d --name <tmp_container_name> -v <volume_name>:/path/inside/container alpine
docker cp /tmp/backup/<volume_name> <tmp_container_name>:/path/inside/container
docker stop <tmp_container_name>
docker rm <tmp_container_name>

# restart the containers
docker restart <container_name>
```

## Schedule Backups

Snapshots and remote syncing is scheduled daily, the full backup archive creation weekly. The default cron configuration can be changed by editing the `/etc/restic/restic.cron` file and restarting the `docker-restic` container with `docker restart <container_name>` or providing a custom cron file using bind mount.

```yml
docker-restic:
  volumes:
    - /path/to/custom.cron:/etc/restic/restic.cron
```

## Available Commands

The following commands are available inside the `docker-restic` container:

- `backup`:
  - stop containers
  - create snapshot
  - start containers
  - prune snapshots
  - check integrity
- `extract`:
  - create archive
  - prune archives
  - check integrity
- `sync`:
  - sync remote
  - check integrity
- `check`:
  - check snapshot integrity
  - check archive integrity
  - check remote integrity

Additionally all official Restic CLI commands are available. [Restic Man Page](https://www.mankier.com/1/restic)

## Contributing

We welcome contributions to Docker-Restic! If you have suggestions, bug reports, or would like to contribute new features, please feel free to submit a pull request or open an issue on the GitHub repository.
