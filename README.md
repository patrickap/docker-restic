# Docker-Restic

Docker-Restic is a lightweight wrapper designed to streamline the use of Restic, particularly for container backups. By parsing a configuration file, Docker-Restic exposes specified commands through the command-line interface (CLI).

## Key Features

- **User-Friendly CLI**: Offers a robust and intuitive command-line interface.
- **Restic Integration**: Supports all available Restic commands, arguments, options and more.
- **Multiple Repositories**: Enables seamless management of multiple repositories and backup locations.
- **Centralized Configuration**: Utilizes a central configuration file for all custom commands.
- **Custom Hooks**: Allows the definition of hooks to execute tailored workflows.
- **Custom Commands**: Facilitates the creation of custom commands for maximum flexibility.
- **Automation Capabilities**: Supports the scheduling of commands for automated backup operations.
- **Non-root Container**: Operates as a non-root container by default, adhering to best security practices.
- **Optional Capabilities**: Offers optional capabilities to read data from different owners if necessary.

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
  # -v $(pwd)/docker-restic.yml:/srv/restic/docker-restic.yml:ro \
  # -v $(pwd)/docker-restic.cron:/srv/restic/docker-restic.cron:ro \

  # Back up the named volume "data"
  -v data:/source/data:ro \
  # Bind mount the backups to the host
  -v ~/backups:/target \
  -v restic-config:/srv/restic \
  -v /etc/localtime:/etc/localtime:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --secret restic-password \
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
      # - ./docker-restic.yml:/srv/restic/docker-restic.yml:ro
      # - ./docker-restic.cron:/srv/restic/docker-restic.cron:ro

      # Back up the named volume "data"
      - data:/source/data:ro
      # Bind mount the backups to the host
      - ~/backups:/target
      - restic-config:/srv/restic
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    secrets:
      - restic-password

volumes:
  restic-config:
  data:
    external: true

secrets:
  restic-password:
    file: /path/to/restic-password.txt
```

**Notes:**

- For security reasons, it is recommended to mount external volumes for backup as read-only using `:ro`.
- Ensure to bind mount your container backups to a custom location on the host for accessibility.
- The `DAC_READ_SEARCH` capability might be required when backing up multiple volumes with different owners or restricted permissions. This capability allows Docker-Restic to read all directories.

2. **Configure the Docker-Restic Container**

Docker-Restic provides default configurations to help you get started quickly. The following commands are supported out of the box:

- `init`: Initializes a repository at `/target/repository` and expects the password file at `/run/secrets/restic-password`. This must be called once manually.
- `backup`: Stops all necessary containers and creates a snapshot of data mounted at `/source`. On successful execution, it automatically calls `forget`, `check`, and restarts the containers.
- `forget`: Prunes old backup snapshots based on the specified policy.
- `check`: Checks the integrity of the repository.
- `container-start`: Starts all containers labeled `docker-restic.container.stop=true`.
- `container-stop`: Stops all containers labeled `docker-restic.container.stop=true`.

The entire backup process is scheduled once a day at 00:00. If this is not sufficient, the configurations can be modified or overwritten completely. Bind mount your custom configurations like this:

- `docker-restic.yml`: `/srv/restic/docker-restic.yml`
- `docker-restic.cron`: `/srv/restic/docker-restic.cron`

Do not forget to restart the container.

5. **Configure rclone (Optional)**

Remote syncing of backups can be configured with `rclone`. This can be done either by bind mounting the `rclone.conf` to `/srv/restic/rclone.conf` or by running `rclone config` inside the `docker-restic` container. Restic itself supports rclone as a backend. Alternatively, it's possible to run rclone via the Docker-Restic CLI using custom commands.

## Configuration Reference

```yml
# Anchors are constructs that can be reused throughout the config.
# This is especially useful for defining Restic repositories.
# It's not related to Docker Restic.
repository: &repository
  repo: "/target/repository"
  password-file: "/run/secrets/restic-password"

commands:
  # Specify the command name which can be run using the `docker-restic` cli
  # This command config is equivalent to: restic backup /source --repo /srv/restic/repository --password-file /run/secrets/password --tag snapshot --verbose --exclude *.secret --exclude *.bin --exclude-larger-than 2048
  backup:
    # Specify the command to run
    command: ["restic", "backup", "/source"]
    # Alternative syntax
    # command:
    #   - restic
    #   - backup
    #   - /source
    options:
      # Maps directly to command line options
      # Can be either of type boolean, string, integer, or list
      # Anchor aliases can easily be used to reuse common options
      <<: *repository
      tag: snapshot
      verbose: true
      # Every option can be specified with prefix if needed
      # Defaults to "--" (e.g. --verbose)
      # --verbose: true
      # -verbose: true
      exclude:
        - "*.secret"
        - "*.bin"
      exclude-larger-than: 2048
    hooks:
      # Runs before
      pre:
        - <command_name>
      # Runs after
      post:
        # Multiple hooks are supported
        # If a hook fails, the following commands are skipped
        - <command_name>
        - <command_name>
      # Runs only on success
      success:
        - <command_name>
      # Runs only on failure
      failure:
        - <command_name>
```

The configured command named `backup` can now be executed using the `docker-restic` CLI:

```bash
docker-restic run backup
```

## Advanced Configuration

Custom commands besides Restic can easily be added to the config:

```yml
commands:
  sync:
    command: ["rclone", "sync", "/from", "to:remote"]
```

It’s also possible to execute complex shell commands that require interpretation by a specific shell like `/bin/sh -c`.

```yml
commands:
  id:
    command: ["/bin/sh", "-c", "echo $(id)"]
```

It's also possible to change the position of the applied command options. By default, they get added at the end of the command automatically. To change it, run the command in a new shell process and access them using the special variable `$@`. Use `--` to signal the end of options for the `/bin/sh` command. Any arguments after `--` will be treated as positional parameters or arguments for the command string executed by `/bin/sh -c`. Try it in the terminal:

```bash
/bin/sh -c 'echo ${@}' -- --option-1 --option-2 --option-3
```

In the example below, during execution `${@}` gets replaced with the actual options.

```yml
commands:
  hello-world:
    command: ["/bin/sh", "-c", "echo ${@}; echo 'world!'", "--"]
    options:
      hello: true
```

To run commands repeatedly or group commands to a workflow, hooks are suitable. This would run the same command three times:

```yml
commands:
  hello-world:
    command: ["/bin/sh", "-c", "echo 'hello world!'"]
    hooks:
      post:
        - hello-world
        - hello-world
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

We welcome contributions to Docker-Restic! If you have suggestions, bug reports, or would like to contribute new features, please feel free to submit a pull request or open an issue on the GitHub repository.
