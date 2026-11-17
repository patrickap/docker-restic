# Docker-Restic

Docker-Restic is a lightweight wrapper around Restic particularly designed to be used for container backups.

## Key Features

- **CLI**: Simple command-line interface built on Just modules (`exec`, `profile`, `container`).
- **Profiles**: Restic flags are resolved per profile from a single `profiles.toml`.
- **Customization**: Supports custom commands.
- **Automation**: Schedules commands for backups.
- **Container control**: Stops containers during a backup and runs commands inside containers via labels.
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
  # -v $(pwd)/profiles.toml:/etc/docker-restic/profiles.toml:ro \
  # -v $(pwd)/jobs.cron:/etc/docker-restic/jobs.cron:ro \

  # Persist the restic backup data in a named volume
  -v docker-restic-data:/var/lib/docker-restic \

  # Optional: Bind mount the restic backups to your host system
  # -v ~/docker-restic-data:/var/lib/docker-restic \

  # Example: Include any container volume in your backup
  # -v app-data:/app/data:ro \

  # The password file of the profile, here the profile `default`
  -v $(pwd)/restic-password:/run/secrets/default/restic-password:ro \
  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
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

      # Persist the restic backup data in a named volume
      - docker-restic-data:/var/lib/docker-restic

      # Optional: Bind mount the restic backups to your host system
      # - ~/docker-restic-data:/var/lib/docker-restic

      # Example: Include any container volume in your backup
      # - app-data:/app/data:ro

      # The password file of the profile, here the profile `default`
      - ./restic-password:/run/secrets/default/restic-password:ro
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro

volumes:
  docker-restic-data:
  app-data:
    external: true
```

**Notes:**

- For security reasons, it is recommended to mount external volumes for backup as read-only using `:ro`.
- Ensure to bind mount your container backups to a custom location on the host for accessibility.
- The `DAC_READ_SEARCH` capability might be required when backing up multiple volumes with different owners or restricted permissions. This capability allows Docker-Restic to read all directories.
- The container paths are fixed. Change the mounts on the host side, not the paths inside the container.

## Volumes and Environment

| Path | Purpose |
| --- | --- |
| `/var/lib/docker-restic` | State, repositories (`/var/lib/docker-restic/<profile>/repository`). This is also the working directory. |
| `/etc/docker-restic` | Configuration: `profiles.toml` and `jobs.cron`. Can be mounted read-only. |
| `/var/cache/docker-restic` | Restic cache (`RESTIC_CACHE_DIR`). |
| `/run/secrets/<profile>/restic-password` | Password file of a profile. |

| Variable | Default | Description |
| --- | --- | --- |
| `UID` | `1234` | User id of the `restic` user. Ownership of the volumes is adjusted on start if it differs. |
| `GID` | `1234` | Group id of the `restic` user. |

The repository of a profile is initialized automatically on start if its password file exists and the repository does not.

2. **Configure the Docker-Restic Container**

Docker-Restic provides default configurations to help you get started quickly. Optionally it's possible to mount custom configurations. During first launch a Restic repository will be created for each profile at `/var/lib/docker-restic/<profile>/repository`.

You can add additional master keys to a Restic repository using `restic key add`. This lets you grant access without sharing the primary key. Make sure to run this command as the intended user, since the key will be owned by that user. Separate keys provide the same access level but can be revoked individually. They also protect the main key and reduce the risk of its exposure.

The entire backup process is scheduled once a day at 00:00. Depending on your requirements you will need to provide your own configurations or modify the existing ones.

- `profiles.toml`: `/etc/docker-restic/profiles.toml`
- `jobs.cron`: `/etc/docker-restic/jobs.cron`

Do not forget to restart the container.

## Profiles

A profile describes how Restic is called for one repository. All profiles live in `profiles.toml`:

```toml
[default]
repo = "${DOCKER_RESTIC_DATA_DIR}/default/repository"
password-file = "/run/secrets/default/restic-password"

[default.forget]
keep-daily = 7
keep-weekly = 4
prune = true

[default.check]
read-data = true
```

- A top-level table (`[default]`) is a profile. Its scalar keys are passed as `--flags` to every Restic command, for example `--repo ...`.
- A nested table (`[default.forget]`) adds flags for that command only. The arguments are matched one after another as a path below the profile: `forget check` looks for `[default.forget.check]`, not `[default.check]`. Words starting with `-` do not change the path.
- A value of `true` becomes a bare flag (`--prune`), `false` is omitted.
- `${...}` is replaced by environment variables, e.g. `${DOCKER_RESTIC_DATA_DIR}`.
- Add as many profiles as needed. Each profile needs its own password file at `/run/secrets/<profile>/restic-password`.

## Configuration Reference

Docker-Restic utilizes Just under the hood, which is a powerful command runner. Make sure to checkout the [documentation](https://just.systems/man/en) on how to configure it. The commands are executed using the `docker-restic` alias and are grouped into modules:

```bash
docker-restic                 # list the available commands
docker-restic profile         # list the profile commands
docker-restic container       # list the container commands
```

| Command | Description |
| --- | --- |
| `docker-restic exec --profile <profile> <restic-args>` | Runs Restic with the flags of the profile, e.g. `exec --profile default snapshots`. |
| `docker-restic exec --profile <profile> --suspend <restic-args>` | Same, but stops the labeled containers before and starts them again afterwards. |
| `docker-restic profile list` | Lists all profiles. |
| `docker-restic profile expand <profile> <restic-args>` | Prints the resolved Restic arguments of a profile. |
| `docker-restic profile get <profile> <key>` | Prints a single value of a profile. |
| `docker-restic container start` / `stop` | Starts or stops all containers labeled `docker-restic.container.stop=true`. |
| `docker-restic container exec` | Runs the command of the label `docker-restic.container.exec` inside the labeled containers. |

## Scheduling

Jobs are scheduled with [supercronic](https://github.com/aptible/supercronic) using `/etc/docker-restic/jobs.cron`. By default every profile is backed up daily at 00:00:

```cron
0 0 * * * for profile in $(docker-restic profile list); do docker-restic exec --profile "${profile}" --suspend backup && docker-restic exec --profile "${profile}" forget && docker-restic exec --profile "${profile}" check; done
```

## Container Labels

Containers are controlled through labels:

```yml
labels:
  # Stopped while a backup with --suspend is running
  - docker-restic.container.stop=true
  # Command executed by `docker-restic container exec`
  - docker-restic.container.exec=<command>
```

## Manual Backups

For manual backups, simply connect to the container. Run `docker-restic` to list all available commands. The `docker-restic` command runs as the user inside the container (by default `restic`) automatically, which prevents the container from writing files as root which the non-root user can't access afterwards. If it happened per accident, run `chown -R restic:restic <directory>` to fix the permissions.

```bash
docker exec -it <container_name> /bin/sh
docker-restic exec --profile default snapshots
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
