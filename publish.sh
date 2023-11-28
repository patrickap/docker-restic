#!/bin/sh

image_name="patrickap/docker-restic"

# Prompt user for version tag
read -p "Enter the version tag for the Docker image: " version_tag

# Log user in to Docker Hub
docker login

# Build the Docker image
docker build . -t "${image_name}:${version_tag}"

# Push the image to Docker Hub
docker push "${image_name}:${version_tag}"

# Remove the locally built image
docker rmi "${image_name}:${version_tag}"
