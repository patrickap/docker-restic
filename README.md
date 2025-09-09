# Docker-Restic

Docker-Restic is a lightweight wrapper around Restic particularly designed to be used for container backups.

## Key Features

- **CLI**: Simple command-line interface.
- **Configuration**: Reads settings from files.
- **Customization**: Supports custom commands.
- **Automation**: Schedules commands for backups.
- **Non-root**: Runs as a non-root container.
- **Capabilities**: Can read volumes owned by different users.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull the Docker-Restic image from Docker Hub and run the container with the specified configurations:

```bash
docker run -d \
  --name docker-restic \
  --restart always \

  # Optional: Add capabilities to read directories of different owners
  # --cap-add DAC_READ_SEARCH \

  # Optional: Mount custom configurations
  # -v $(pwd)/docker-restic.conf:/etc/docker-restic/docker-restic.conf:ro \
  # -v $(pwd)/docker-restic.cron:/etc/docker-restic/docker-restic.cron:ro \
  # -v $(pwd)/rclone.conf:/etc/docker-restic/rclone.conf:ro \

  # Persist the restic backup data in a named volume
  -v docker-restic-data:/var/lib/docker-restic \

  # Optional: Bind mount the restic backups to your host system
  # -v ~/docker-restic-data:/var/lib/docker-restic \

  # Example: Include any container volume in your backup
  # -v app-data:/app/data:ro \

  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --secret restic-password \
  --secret rclone-password \
  patrickap/docker-restic:latest
```

Alternatively, you can use Docker Compose:

```yml
services:
  docker-restic:
    image: patrickap/docker-restic:latest
    restart: always
    # Optional: Add capabilities to read directories of different owners
    # cap_add:
    # - DAC_READ_SEARCH
    volumes:
      # Optional: Mount custom configurations
      # - ./docker-restic.conf:/etc/docker-restic/docker-restic.conf:ro
      # - ./docker-restic.cron:/etc/docker-restic/docker-restic.cron:ro
      # - ./rclone.conf:/etc/docker-restic/rclone.conf:ro

      # Persist the restic backup data in a named volume
      - docker-restic-data:/var/lib/docker-restic

      # Optional: Bind mount the restic backups to your host system
      # - ~/docker-restic-data:/var/lib/docker-restic

      # Example: Include any container volume in your backup
      # - app-data:/app/data:ro

      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    secrets:
      - restic-password
      - rclone-password

volumes:
  restic-data:
  app-data:
    external: true

secrets:
  restic-password:
    file: /run/secrets/restic-password.txt
  rclone-password:
    file: /run/secrets/rclone-password.txt
```

**Notes:**

- For security reasons, it is recommended to mount external volumes for backup as read-only using `:ro`.
- Ensure to bind mount your container backups to a custom location on the host for accessibility.
- The `DAC_READ_SEARCH` capability might be required when backing up multiple volumes with different owners or restricted permissions. This capability allows Docker-Restic to read all directories.

2. **Configure the Docker-Restic Container**

Docker-Restic provides default configurations to help you get started quickly. Optionally it's possible to mount custom configurations. During first launch a Restic repository will be created for you at `/var/lib/docker-restic/repository` as well as an encrypted Rclone configuration at `/etc/docker-restic/rclone.conf`.

You can add additional master keys to a Restic repository using `restic key add`. This lets you grant access without sharing the primary key. Separate keys provide the same access level but can be revoked individually. They also protect the main key and reduce the risk of its exposure.

To add a new remote backup target either run `rclone config` inside the container or mount an existing encrypted configuration. By default all Rclone remotes with names ending in `_default` will automatically be used as backup targets and kept in sync. Since the remote path cannot be known in advance you need to specify it explicitly in your configuration. This can be achieved using aliases. In this example your backups will be synced to the `drive` remote at the path `/backup`.

```ini
[drive]
type = drive
client_id = <client_id>
client_secret = <client_secret>
token = <token>

[drive_default]
type = alias
remote = drive:/backup
```

The entire backup process is scheduled once a day at 00:00. Depending on your requirements you will need to provide your own configurations or modify the existing ones.

- `docker-restic.conf`: `/etc/docker-restic/docker-restic.conf`
- `docker-restic.cron`: `/etc/docker-restic/docker-restic.cron`
- `rclone.conf`: `/etc/docker-restic/rclone.conf`

Do not forget to restart the container.

## Configuration Reference

Docker-Restic utilizes Just under the hood, which is a powerful command runner. Make sure to checkout the [documentation](https://just.systems/man/en) on how to configure it. The configured commands should be executed using the `docker-restic` alias:

```bash
docker-restic <command-name>
```

## Manual Backups

For manual backups, simply connect to the container. A lot of useful commands are provided by default. Run `docker-restic -l` to list all available commands. It's important to run the commands as the user inside the container (by default `restic`) to prevent the container from writing files as root which the non-root user can't access afterwards. If it happened per accident, run `chown -R restic:restic <directory>` to fix the permissions:

```bash
docker exec -u <user> -it <container_name> /bin/sh
```

## Restore from Backup

To restore a backup, a new Docker volume with the correct name must be created including the contents of the backup. After restarting the containers, the data should be mounted and restored:

```bash
# check restic repository
restic -r /path/to/repository check --read-data

# dump restic backup
restic -r /path/to/repository dump latest / > backup.tar

# untar the backup
tar -xvf backup.tar -C /tmp/backup

# stop the containers
docker stop <container_name>

# use a temporary container to create the volume and copy the backup
docker volume create <volume_name>
docker run --rm -it -v <volume_name>:/to -v <path_to_backup>:/from alpine /bin/sh -c 'cp -av /from/. /to'

# restart the containers
docker restart <container_name>
```

**Warning:**
If you're using Google Drive they may add back file extensions to encrypted files during the download or compression process which can result in a corrupted `restic` repository. To avoid this ensure to remove any added extensions inside the `repository/data` directory. An example of this would be a file at `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b` being downloaded as `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b.mp3`. In this example it is required to remove the automatically added extension `.mp3` to avoid repository corruption and be able to read the backup.

## Contributing

To run Docker-Restic locally, you have two options: either build the Docker image from the provided Dockerfile and execute it, or use `docker compose`. To publish a release, use the command `just release <patch|minor|major>`. This command will automatically increment the semantic version accordingly.
