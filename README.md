# Docker-Restic

Docker-Restic is a lightweight wrapper designed to streamline the use of Restic, particularly for container backups. By parsing a configuration file, Docker-Restic exposes specified commands through the command-line interface (CLI).

## Key Features

- **User-Friendly CLI**: Offers a robust and intuitive command-line interface.
- **Centralized Configuration**: Utilizes a central configuration file for all custom commands.
- **Custom Commands**: Facilitates the creation of custom commands for maximum flexibility.
- **Automation Capabilities**: Supports the scheduling of commands for automated backup operations.
- **Non-root Container**: Operates as a non-root container by default, adhering to best security practices.
- **Optional Capabilities**: Offers optional capabilities to read volumes from different owners if necessary.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull the Docker-Restic image from the official Docker Hub repository and run the container with the specified configurations:

```bash
docker run -d \
  --name docker-restic \
  --restart always \

  # Optional: Add capabilities to read directories of different owners
  # --cap-add DAC_READ_SEARCH \

  # Optional: Overwrite the default configuration
  # -v $(pwd)/restic.conf:/srv/restic/config/restic.conf:ro \
  # -v $(pwd)/restic.cron:/srv/restic/config/restic.cron:ro \

  # Back up the named volume "data"
  -v data:/source/data:ro \
  # Bind mount the backups to the host
  -v ~/backups:/target \
  -v restic-config:/srv/restic \
  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --secret restic-password \
  --secret rclone-password \
  patrickap/docker-restic:latest
```

Alternatively, you can use Docker Compose:

```yml
version: "3.7"

services:
  docker-restic:
    image: patrickap/docker-restic:latest
    restart: always
    # Optional: Add capabilities to read directories of different owners
    # cap_add:
    # - DAC_READ_SEARCH
    volumes:
      # Optional: Overwrite the default configuration
      # - ./restic.conf:/srv/restic/config/restic.conf:ro
      # - ./restic.cron:/srv/restic/config/restic.cron:ro

      # Back up the named volume "data"
      - data:/source/data:ro
      # Bind mount the backups to the host
      - ~/backups:/target
      - restic-config:/srv/restic
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    secrets:
      - restic-password
      - rclone-password

volumes:
  restic-config:
  data:
    external: true

secrets:
  restic-password:
    file: /path/to/restic-password.txt
  rclone-password:
    file: /path/to/rclone-password.txt
```

**Notes:**

- For security reasons, it is recommended to mount external volumes for backup as read-only using `:ro`.
- Ensure to bind mount your container backups to a custom location on the host for accessibility.
- The `DAC_READ_SEARCH` capability might be required when backing up multiple volumes with different owners or restricted permissions. This capability allows Docker-Restic to read all directories.

2. **Configure the Docker-Restic Container**

Docker-Restic provides default configurations to help you get started quickly. A lot of commands are supported out of the box. Run `docker-restic -l` to list all available commands. Run `docker-restic init` once manually. This will initialize a restic repository at `/srv/restic/data/repository` and creates an encrypted rclone configuration at `/srv/restic/config/rclone.conf`.

The entire backup process is scheduled once a day at 00:00. If this is not sufficient, the configurations can be modified or overwritten completely. Bind mount your custom configurations like this:

- `restic.conf`: `/srv/restic/config/restic.conf`
- `restic.cron`: `/srv/restic/config/restic.cron`

Do not forget to restart the container.

## Configuration Reference

Docker-Restic utilizes Just under the hood, which is a powerful command runner. Make sure to checkout the [documentation](https://just.systems/man/en) on how to configure it. The configured commands should be executed using the `docker-restic` alias:

```bash
docker-restic <command-name>
```

## Manual Backups

For manual backups, simply connect to the container. It's important to run the container as the user inside the container (by default `restic`) to prevent the container from writing files as root which the non-root user can't access afterwards. If it happened per accident, run `chown -R restic:restic <directory>` to fix the permissions:

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
