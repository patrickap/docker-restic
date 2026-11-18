# Docker-Restic

Docker-Restic is a lightweight wrapper around Restic, specifically designed for container backups.

## Key Features

- **CLI**: Simple and intuitive command-line interface.
- **Profiles**: Flexible profile system for managing multiple Restic repositories.
- **Automation**: Built-in cron scheduling for recurring backup jobs.
- **Safety**: Automatic container management during backups via service labels.
- **Non-root**: Secure operation by running as a non-root user with minimal privileges.

## Getting Started

To get started with Docker-Restic, deploy the container using Docker or Docker Compose.

### Docker

```bash
docker run -d \
  --name docker-restic \
  --restart always \
  # Optional: Add capabilities to read directories owned by other users
  # --cap-add DAC_READ_SEARCH \

  # Optional: Mount custom configuration files
  # -v ./profiles.toml:/etc/docker-restic/profiles.toml:ro \
  # -v ./jobs.cron:/etc/docker-restic/jobs.cron:ro \

  # Persist Restic data in a named volume
  -v data:/var/lib/docker-restic \

  # Optional: Bind mount Restic data to the host
  # -v ./data:/var/lib/docker-restic \

  # Example: Include a container volume in the backup
  # -v app-data:/app/data:ro \

  # Mount the password file for the default profile
  -v ./restic-password:/run/secrets/default/restic-password:ro \

  # Synchronize the container's local time with the host
  -v /etc/localtime:/etc/localtime:ro \

  # Optional: Mount the Docker socket for container management (start|stop|exec)
  # -v /var/run/docker.sock:/var/run/docker.sock:ro \
  patrickap/docker-restic:latest
```

### Docker Compose

Alternatively, deploy Docker-Restic using Docker Compose:

```yaml
services:
  docker-restic:
    image: patrickap/docker-restic:latest
    restart: always
    # Optional: Add capabilities to read directories owned by other users
    # cap_add:
    #   - DAC_READ_SEARCH
    volumes:
      # Optional: Mount custom configuration files
      # - ./profiles.toml:/etc/docker-restic/profiles.toml:ro
      # - ./jobs.cron:/etc/docker-restic/jobs.cron:ro

      # Persist Restic data in a named volume
      - data:/var/lib/docker-restic

      # Optional: Bind mount Restic data to the host
      # - ./data:/var/lib/docker-restic

      # Example: Include a container volume in the backup
      # - app-data:/app/data:ro

      # Mount the password file for the default profile
      - ./restic-password:/run/secrets/default/restic-password:ro

      # Synchronize the container's local time with the host
      - /etc/localtime:/etc/localtime:ro

      # Optional: Mount the Docker socket for container management (start|stop|exec)
      # - /var/run/docker.sock:/var/run/docker.sock:ro
```

### Notes

- Bind mount Docker-Restic volumes to a dedicated location to ensure data persists across container recreations.
- The `DAC_READ_SEARCH` capability may be required when backing up volumes owned by different users or with restricted permissions. It allows Docker-Restic to read directories that would otherwise be inaccessible.
- For security and safety reasons, external volumes included in backups should be mounted as read-only using `:ro`.
- To enable container management during backups, mount the Docker socket. This allows Docker-Restic to stop and start containers and execute custom backup commands.

## Configuration

Docker-Restic provides default configurations to simplify the initial setup. On startup, the container:

1. Checks and corrects potential permission issues.
2. Initializes a Restic repository for each configured profile if it does not already exist.
3. Starts the cron process to execute scheduled jobs.

The following configuration files are provided with default configurations:

| File | Location |
| --- | --- |
| `profiles.toml` | `/etc/docker-restic/profiles.toml` |
| `jobs.cron` | `/etc/docker-restic/jobs.cron` |

Custom configuration files can optionally be mounted into the container. After changing the configuration, restart the container to apply the changes.

### Environment Variables

The following environment variables can be set:

| Variable | Default | Description |
| --- | --- | --- |
| `UID` | `1234` | User ID of the `restic` user. |
| `GID` | `1234` | Group ID of the `restic` user. |
| `DOCKER_RESTIC_DATA_DIR` | `/var/lib/docker-restic` | Directory for application state and data. |
| `DOCKER_RESTIC_CONFIG_DIR` | `/etc/docker-restic` | Directory containing configuration files, such as `profiles.toml` and `jobs.cron`. |
| `RESTIC_CACHE_DIR` | `/var/cache/restic` | Directory used by Restic to store its cache. |

### Profiles

A profile defines how Restic is invoked for a specific repository. All profiles are configured in `profiles.toml`. All Restic-specific flags are supported. A `default` profile is preconfigured out of the box.

Example:

```toml
[default]
repo = "${DOCKER_RESTIC_DATA_DIR}/default/repository"
password-file = "/run/secrets/default/restic-password"

[default.forget]
keep-daily = 7
...
```

The following rules apply to profile configuration:

- A top-level table (e.g. `[default]`) defines a profile. Its scalar values are passed as flags to every Restic command.
- A nested table (e.g. `[default.forget]`) defines additional flags that are passed only to the corresponding Restic command.
- Arguments are matched sequentially against the configured paths.
- Custom flags and commands inserted between arguments are preserved and do not affect path matching.
- All values are converted to strings. A value of `false` is omitted, while `true` is retained as a boolean flag.
- Environment variables referenced in the configuration file are automatically substituted.
- Any number of profiles can be configured. Each profile requires its own configuration.

### Scheduling

Backup jobs are defined in `jobs.cron` and executed automatically according to their cron schedules. By default, every configured profile is backed up daily at `00:00`.

Example:

```cron
0 0 * * * for profile in $(docker-restic profile list); do docker-restic exec --profile "${profile}" --container-stop backup ...; done
```

### Container Labels

Docker-Restic uses Docker labels to control container behavior during backups and execute commands inside containers.

The following labels are supported:

```yaml
labels:
  # Stop the container during a backup and restart it afterwards
  - docker-restic.container.stop=true

  # Execute a custom command inside the container during a backup
  - docker-restic.container.exec=<command>
```

### Master Keys

Additional master keys can be added to a Restic repository using:

```bash
docker-restic exec --profile <profile> key add
```

This allows you to grant repository access without sharing the primary key.

Make sure to execute this command as the intended user, as the resulting key will be owned by that user.

Each key provides the same level of access and can be revoked individually. Using separate keys also helps protect the primary key and reduces the risk of exposing it.

## CLI Reference

All Restic-related commands should be executed through the `docker-restic` CLI.

To list the available commands:

```bash
docker-restic --list
```

### Commands

| Command | Description |
| --- | --- |
| `docker-restic exec --profile <profile> <restic_args>` | Executes Restic using the configuration of the specified profile. |
| `docker-restic profile list` | Lists all configured profiles. |
| `docker-restic profile expand <profile> <restic_args>` | Prints the resolved Restic arguments for the specified profile. |
| `docker-restic profile get <profile> <path>` | Prints a single configuration value from the specified profile. |
| `docker-restic container start` / `stop` | Starts or stops all containers labeled `docker-restic.container.stop=true`. |
| `docker-restic container exec` | Executes the command defined by the `docker-restic.container.exec` label inside labeled containers. |

## Manual Backups and Dumps

For manual backups or other operations, connect to the Docker-Restic container using a shell:

```bash
docker exec -it <container_name> /bin/sh
```

You can then execute the available `docker-restic` commands directly.

## Restore from Backup

To restore a backup, create a Docker volume with the appropriate name and copy the restored data into it.

The following example demonstrates how to restore a backup:

```bash
# List all snapshots
docker-restic exec --profile <profile> snapshots

# Check the Restic repository
docker-restic exec --profile <profile> check --read-data

# Export the latest backup
docker-restic exec --profile <profile> dump latest / > backup.tar

# Extract the backup
tar -xvf backup.tar -C /tmp/backup

# Stop the containers
docker stop <container_name>

# Create a Docker volume
docker volume create <volume_name>

# Copy the restored data into the volume using a temporary container
docker run --rm -it \
  -v <volume_name>:/to \
  -v <path_to_backup>:/from \
  alpine /bin/sh -c 'cp -av /from/. /to'

# Restart the containers
docker restart <container_name>
```

## Contributing

To run Docker-Restic locally, you can either build and run the Docker image using the provided Dockerfile or use Docker Compose.

To publish a new release, run:

```bash
just release <patch|minor|major>
```

This command automatically increments the semantic version according to the specified release type.
