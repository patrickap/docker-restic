# docker-restic

Docker-Restic is a small wrapper that simplifies the use of restic especially for container backups. It parses a configuration file and makes the specified commands available via the CLI.

## Features

- **Simple CLI**: Provides a robust command-line interface.
- **Restic Commands**: Supports all available restic commands, arguments, and options.
- **Multiple Repositories** Supports multiple repositories and backup locations.
- **Config File**: Utilizes a central configuration file for all custom commands.
- **Custom Hooks**: Allows defining hooks to run custom workflows.
- **Custom Commands**: Allows defining custum commands for max flexibility.
- **Automation**: Supports scheduling of commands out of the box.
- **Non-root Container**: Runs as a non-root container by default.
- **Capabilities**: Optional capabilities to read data from other users.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull Docker-Restic from the official Docker Hub repository and run the container image:

```bash
docker run -d \
  --name docker-restic \
  --restart always \
  --cap-add DAC_READ_SEARCH \
  -v $(pwd)/docker-restic.yml:/srv/docker-restic/docker-restic.yml:ro \
  -v $(pwd)/docker-restic.cron:/srv/docker-restic/docker-restic.cron:ro \
  -v docker-restic-data:/srv/docker-restic \
  -v media:/media:ro \
  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --secret restic-password \
  patrickap/docker-restic:latest
```

**Note:**

- For safety reasons it's recommended to mount external volumes to backup as read-only using `:ro`.
- Make sure to bind mount your container backups to a custom location on the host to be able to access them at any time.
- It may be necessary to add the `DAC_READ_SEARCH` capability to the container when backing up multiple volumes from different owners or with restricted permissions. This capability will allow Docker-Restic to read all directories.

Alternatively with docker compose:

```yml
version: "3.7"

services:
  docker-restic:
    image: patrickap/docker-restic:latest
    restart: always
    cap_add:
      - DAC_READ_SEARCH
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

2. Configure the Docker-Restic container:

<!-- TODO: mention and create default config.yml -->

Docker-Restic can be configured using a central `docker-restic.yml` file. Create a local config and bind mount it into the container to `/srv/docker-restic/docker-restic.yml`

```yml
# specify a repository
# this is user specific and ignored by Docker-Restic
# however standard yaml anchors can be used to reuse common values
repository: &repository
  repo: "/srv/docker-restic/backup/repository"
  password-file: "/run/secrets/restic-password"

commands:
  # equivalent to: restic backup /media --repo /srv/restic/repository --password-file /run/secrets/password --tag snapshot --verbose --exclude *.secret --exclude *.bin --exclude-larger-than 2048
  snapshot:
    # specify command to run
    command: ["restic", "backup", "/media"]
    # command:
    #   - restic
    #   - backup
    #   - /media
    options:
      # maps directly to restic command options
      # can be either boolean, string, integer or a list
      # anchors can be used to reuse common options
      <<: *repository
      tag: snapshot
      # options default to prefix "--"
      # in this case --verbose
      # can also be added manually with "-" or "--" for compatability
      verbose: true
      # -verbose: true
      # --verbose: true
      exclude:
        - "*.secret"
        - "*.bin"
      exclude-larger-than: 2048
    hooks:
      # runs before
      pre:
        - <command_name>
      # runs after
      post:
        # multiple hooks are supported
        # if a hook fails the remaining commands are skipped
        - <command_name>
        - <command_name>
      # runs only on success
      success:
        - <command_name>
      # runs only on failure
      failure:
        - <command_name>
```

configured command `snapshot` above can later be called using the `docker-restic` CLI.

```bash
docker-restic run snapshot
```

<!-- TODO: mention and create default config.cron -->

To schedule commands create a `docker-restic.cron` file and bind mount it to `/srv/docker-restic/docker-restic.cron`

```bash
# daily backup
0 0 * * * docker-restic run snapshot
```

The command gets scheduled on container startup / restart.

5. Configure rclone (optional)

Remote syncing of backups can be configured with `rclone`. Either by bind mount the `rclone.conf` to `/srv/docker-restic/rclone.conf` or run `rclone config` inside the `docker-restic` container. Restic itself supports rclone as backend. Alternatively it's possible to run rclone via the Docker-Restic CLI using custom commands.

## Advance Config

custom commands (e.g. for rclone) can easily be added to the config yaml

```yml
commands:
  sync:
    command: ["rclone", "sync", "/from", "to:remote"]
```

its also possible to run in a new shell for shell process like `/bin/sh -c` for accessing env variables `$ENV_VAR` or using shell operators like pipes etc.

```yml
commands:
  home:
    command: ["/bin/sh", "-c", "echo $HOME"]
```

its also possible to change the position of the applied command options. by default they get added at the end of the command autoamtically. to change it run the command in new shell process and access them using special variable `$@`. use `--` to tell the shell that this is the end of the /bin/sh -c command. try it in terminal.

```bash
/bin/sh -c "echo ${@}" -- --option-1 --option-2 --option-3
```

in exmaple below during execution `${@}` gets replaced with the actual options. this result in `custom-command --test another-argument`.

```yml
commands:
  custom:
    command: ["/bin/sh", "-c", "custom-command ${@} another-argument", "--"]
    options:
      test: true
```

to run commands multiple times or group commands like for a workflow. hooks are a good fit. this example would run the same command 3 times.

```yml
commands:
  custom:
    command: ["/bin/sh", "-c", "echo $HOME"]
    hooks:
      post:
        - custom
        - custom
```

## Manual Backups

for manual backups simply connect to the container. its important to run container as the user inside the container (by default `restic`) to prevent the container from writing files as root user / owner which the non-root suer cant access afterwards. if it happened per accident run `chown -R restic:restic /dir` to fix the permissions

```bash
docker exec -u <user> -it <container_name> /bin/sh
```

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
