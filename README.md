# docker-restic

Docker-Restic is a small wrapper that simplifies the use of restic especially for container backups. It parses a configuration file and makes the specified commands available via the CLI.

## Features

- **Simple CLI**: Provides a robust command-line interface.
- **Restic Commands**: Supports all available restic commands, arguments, and flags.
- **Multiple Repositories** Supports multiple repositories and backup locations.
- **Config File**: Utilizes a central configuration file for all custom commands.
- **Custom Hooks**: Allows defining hooks to run custom commands.
- **Automation**: Supports scheduling of commands out of the box.
- **Non-root Container**: Runs as a non-root container by default.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull the Docker-Restic image from the official Docker Hub repository:

```shell
docker pull patrickap/docker-restic:latest
```

2. Create a local `docker-restic.yml` file.

```yml
repositories:
  default-repository:
    &default-repository # maps directly to restic command flags
    repo: "/srv/docker-restic/repository"
    password-file: "/run/secrets/restic-password"

commands:
  # equivalent to: restic backup /media --repo /srv/restic/repository --password-file /run/secrets/password --tag snapshot --verbose --exclude *.secret --exclude *.bin --exclude-larger-than 2048
  snapshot:
    arguments:
      # maps directly to restic command arguments
      # order is guaranteed
      - backup
      - /media
    flags:
      # maps directly to restic command flags
      # can be either boolean, string, integer or a list
      # anchors can be used to reuse common flags
      <<: *default-repository
      tag: snapshot
      verbose: true
      exclude:
        - "*.secret"
        - "*.bin"
      exclude-larger-than: 2048
    hooks:
      # runs before
      pre: "echo 'pre'"
      # runs after
      post: "echo 'post'"
      # runs only on success
      success: "echo 'success'"
      # runs only on failure
      failure: "echo 'failure'"
```

The configuration gets directly translated into native restic commands. The configured command `snapshot` above can later be called using `docker-restic`.

```bash
docker-restic run snapshot
```

3. Create a local `docker-restic.cron` file.

```bash
# daily
0 0 * * * docker-restic run snapshot
```

The command gets scheduled on container startup.

4. Run the container image:

```bash
docker run -d \
  --name docker-restic \
  --restart always \
  -v $(pwd)/docker-restic.yml:/srv/docker-restic/docker-restic.yml:ro \
  -v $(pwd)/docker-restic.cron:/srv/docker-restic/docker-restic.cron:ro \
  -v docker-restic-data:/srv/docker-restic \
  -v media:/media:ro \
  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --secret restic-password \
  patrickap/docker-restic:latest
```

Alternatively with docker compose:

```yml
version: "3.7"

services:
  docker-restic:
    image: patrickap/docker-restic:latest
    restart: always
    volumes:
      - ./docker-restic.yml:/srv/docker-restic/docker-restic.yml:ro
      - ./docker-restic.cron:/srv/docker-restic/docker-restic.cron:ro
      - docker-restic-data:/srv/docker-restic
      - media:/media:ro
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    secrets:
      - restic-password

  volumes:
    docker-restic-data:
    media:

  secrets:
    restic-password:
      file: /path/to/restic-password.txt
```

**Note:**
For safety reasons it's recommended to mount external volumes as read-only using `:ro`.

```yml
docker-restic:
  volumes:
    - media:/media:ro
```

Make sure to bind mount the restic repository path to a custom location to be able to access it at any time.

```yml
docker-restic:
  volumes:
    - ./repository:/srv/docker-restic/repository
```

It may be necessary to add the `DAC_READ_SEARCH` capability to the container when backing up multiple volumes from different owners or with restricted permissions. This capability will allow Docker-Restic to read all directories.

```yml
docker-restic:
  cap_add:
    - DAC_READ_SEARCH
```

5. Configure rclone (optional)

Remote syncing of backups can be configured with `rclone`. Either by bind mounting the `rclone.conf` to `/srv/docker-restic/rclone.conf` into the container or run `rclone config` inside the `docker-restic` container. Restic itself supports rclone as backend. Alternatively it's possible to run rclone via hooks.

## Restore from Backup

To restore a backup a new Docker volume with the correct name must be created including the contents of the backup. After restarting the containers the data should be mounted and restored.

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
docker run --rm -it -v <volume_name>:/to -v <path_to_backup>:/from alpine ash -c 'cp -av /from/. /to'

# restart the containers
docker restart <container_name>
```

**Warning:**
If you're using Google Drive they may add back file extensions to encrypted files during the download or compression process which can result in a corrupted `restic` repository. To avoid this ensure to remove any added extensions inside the `repository/data` directory. An example of this would be a file at `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b` being downloaded as `respository/data/3f/3f0e4a8c5b71a0b9c7d38e29a87d5a1b23f69b08a5c06f1d2b539c846ee2a070b.mp3`. In this example it is required to remove the automatically added extension `.mp3` to avoid repository corruption and be able to read the backup.

## Contributing

We welcome contributions to Docker-Restic! If you have suggestions, bug reports, or would like to contribute new features, please feel free to submit a pull request or open an issue on the GitHub repository.
