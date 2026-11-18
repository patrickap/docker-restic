# Docker-Restic

Docker-Restic is a lightweight wrapper around Restic particularly designed to be used for container backups.

## Key Features

- **CLI**: Simple and intuitive command-line interface.
- **Profiles**: Flexible profile system for managing multiple Restic configurations.
- **Automation**: Built-in cron scheduling for automating recurring backup jobs.
- **Safety**: Automatically stops containers before backups using service labels.
- **Non-root**: Runs as a non-root user by default.

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
  # -v ./profiles.toml:/etc/docker-restic/profiles.toml:ro \
  # -v ./jobs.cron:/etc/docker-restic/jobs.cron:ro \

  # Persist the restic data in a named volume
  -v data:/var/lib/docker-restic \

  # Optional: Bind mount the restic data to your host system
  # -v ./data:/var/lib/docker-restic \

  # Example: Include any container volume in your backup
  # -v app-data:/app/data:ro \

  # The password file of the profile, here the profile `default`
  -v ./restic-password:/run/secrets/default/restic-password:ro \

  # Bind mount the local time
  -v /etc/localtime:/etc/localtime:ro \

  # Optional: Bind mount the docker socket for exec commands
  # -v /var/run/docker.sock:/var/run/docker.sock:ro \
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
      # - ./profiles.toml:/etc/docker-restic/profiles.toml:ro
      # - ./jobs.cron:/etc/docker-restic/jobs.cron:ro

      # Persist the restic data in a named volume
      - data:/var/lib/docker-restic

      # Optional: Bind mount the restic backups to your host system
      # - ./data:/var/lib/docker-restic

      # Example: Include any container volume in your backup
      # - app-data:/app/data:ro

      # The password file of the profile, here the profile `default`
      - ./restic-password:/run/secrets/default/restic-password:ro

      # Bind mount the local time
      - /etc/localtime:/etc/localtime:ro

      # Optional: Bind mount the docker socket for exec commands
      - /var/run/docker.sock:/var/run/docker.sock:ro
```

**Notes:**

- For security reasons, it is recommended to mount external volumes for backup as read-only using `:ro`.
- Ensure to bind mount your container data to a custom location for quick access during an emergency.
- The `DAC_READ_SEARCH` capability might be required when backing up multiple volumes with different owners or restricted permissions. This capability allows Docker-Restic to read all directories.

## Configuration

Docker-Restic provides default configurations to help you get started quickly. During first launch the container will check and fix potential permission issues and will intialize new Restic repositories for each configured profile and will start the main cron process.

Optionally it's possible to mount custom configurations.
The entire backup process is scheduled once a day at 00:00. Depending on your requirements you will need to provide your own configurations or modify the existing ones.

- `profiles.toml`: `/etc/docker-restic/profiles.toml`
- `jobs.cron`: `/etc/docker-restic/jobs.cron`

Do not forget to restart the container.

### Environment

| Variable | Default | Description |
| --- | --- | --- |
| `UID` | `1234` | User ID of the `restic` user. |
| `GID` | `1234` | Group ID of the `restic` user. |
| `DOCKER_RESTIC_DATA_DIR` | `/var/lib/docker-restic` | Directory for the application state and data. |
| `DOCKER_RESTIC_CONFIG_DIR` | `/etc/docker-restic` | Directory containing configuration files, such as `profiles.toml` and `jobs.cron`. |
| `RESTIC_CACHE_DIR` | `/var/cache/restic` | Directory used by Restic to store its cache. |

The repository of a profile is initialized automatically on start if the repository does not.

### Profiles

A profile describes how Restic is called for one repository. All profiles live in `profiles.toml`:

```toml
[default]
repo = "${DOCKER_RESTIC_DATA_DIR}/default/repository"
password-file = "/run/secrets/default/restic-password"

[default.forget]
keep-daily = 7
prune = true

[default.check]
read-data = true
...
```

- A top-level table (`[default]`) is a profile. Its scalar keys are passed as flags to every Restic command.
- A nested table (`[default.forget]`) adds flags for that command only. The arguments are matched one after another as a path
- Custom flags and commands in-between added during the call are kept in place and do not affect the path matching.
- Values are transformed to strings. A value of `false` is omitted and a value of `true` is kept as a boolean flag
- Environment variabels inside the config file are supported and will get substituted.
- Supported are as many profiles as needed. Each profile needs its own configuration.

### Scheduling

Jobs are scheduled automatically. By default every profile is backed up daily at 00:00:

```cron
0 0 * * * for profile in $(docker-restic profile list); do docker-restic exec --profile "${profile}" --suspend backup ...; done
```

### Container Labels

Containers are controlled through labels:

```yml
labels:
  # Stopped while a backup with --suspend is running
  - docker-restic.container.stop=true
  # Command executed by `docker-restic container exec`
  - docker-restic.container.exec=<command>
```

### Master Keys

You can add additional master keys to a Restic repository using `docker-restic exec --profile <profile> key add`. This lets you grant access without sharing the primary key. Make sure to run this command as the intended user, since the key will be owned by that user. Separate keys provide the same access level but can be revoked individually. They also protect the main key and reduce the risk of its exposure.

## CLI Reference

All Restic related commands should be executed using the `docker-restic` alias:

```bash
docker-restic --list                    # list the available commands
docker-restic profile --list            # list the profile commands
docker-restic container --list          # list the container commands
```

| Command | Description |
| --- | --- |
| `docker-restic exec --profile <profile> <restic_args>` | Runs Restic with the flags of the profile, e.g. `exec --profile default snapshots`. |
| `docker-restic exec --profile <profile> --suspend <restic_args>` | Same, but stops the labeled containers before and starts them again afterwards. |
| `docker-restic profile list` | Lists all profiles. |
| `docker-restic profile expand <profile> <restic_args>` | Prints the resolved Restic arguments of a profile. |
| `docker-restic profile get <profile> <path>` | Prints a single value of a profile. |
| `docker-restic container start` / `stop` | Starts or stops all containers labeled `docker-restic.container.stop=true`. |
| `docker-restic container exec` | Runs the command of the label `docker-restic.container.exec` inside the labeled containers. |

## Manual Backups and Dumps

For manual backups, simply connect to the container.

```bash
docker exec -it <container_name> /bin/sh
```

## Restore from Backup

To restore a backup, a new Docker volume with the correct name must be created including the contents of the backup. After restarting the containers, the data should be mounted and restored:

```bash
# list all snapshots
docker-restic exec --profile <profile> snapshots

# check restic repository
docker-restic exec --profile <profile> check --read-data

# dump restic backup
docker-restic exec --profile <profile> dump latest / > backup.tar

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

## Contributing

To run Docker-Restic locally, you have two options: either build the Docker image from the provided Dockerfile and execute it, or use `docker compose`. To publish a release, use the command `just release <patch|minor|major>`. This command will automatically increment the semantic version accordingly.
