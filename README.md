# Docker-Restic

Docker-Restic is a Docker image that provides an easy way to use restic for container backups. It efficiently handles three key functionalities: incremental backups, full backups and remote synchronization. The backup process is fully automated.

## Features

- **Incremental Backups:** Docker-Restic performs daily snapshots using restic, allowing you to capture changes in your data efficiently.
- **Full Backups:** Docker-Restic automatically exports a weekly tar archive file, providing a full backup of your data.
- **Remote Sync (Optional):** You have the option to enable a remote synchronization using rclone, which ensures that your backups are securely transferred to a remote location.
- **Easy Setup:** All data mounted at `/source` within Docker-Restic is backed up to `/target`. This flexible setup allows you to define the specific directories and volumes you wish to include in your backups.
- **Container Management:** Containers labeled with `restic-stop=true` are gracefully stopped before the backup process and restarted afterward, ensuring data consistency during the backup operation.
- **File Locking:** To prevent concurrent access to backup resources, Docker-Restic utilizes a lockfile mechanism that effectively manages access and avoids conflicts.
- **Flexible Customization:** Docker-Restic offers a high level of customization through various `ARG`s and `ENV`s that can be easily set or overwritten according to your requirements. These customization options provide the flexibility to adapt the backup process to your specific needs.

## Getting Started

To get started with Docker-Restic, follow these steps:

1. Pull the Docker-Restic image from the official Docker Hub repository:

```shell
docker pull patrickap/docker-restic:latest
```

2. Configure the necessary `ARG`s and `ENV`s to suit your backup requirements. Refer to the `Dockerfile` for a complete list of customization options.

3. Run the Docker-Restic container:

```shell
docker run -d --name docker-restic \
    -v /path/to/source:/source \
    -v /path/to/target:/target \
    -e RESTIC_PASSWORD=your-password \
    patrickap/docker-restic:latest
```

4. Monitor the backup process and view the logs:

```shell
docker logs docker-restic
```

## Contributing

We welcome contributions to Docker-Restic! If you have suggestions, bug reports, or would like to contribute new features, please feel free to submit a pull request or open an issue on the GitHub repository.

## License

Docker-Restic is released under the [MIT License](https://opensource.org/licenses/MIT).
