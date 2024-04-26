#!/bin/bash

# Stop all running containers
echo "Stopping all running containers..."
docker stop $(docker ps -aq)

# Remove all containers
echo "Removing all containers..."
docker rm $(docker ps -aq)

# Remove all images
echo "Removing all images..."
docker rmi $(docker images -q)

# Remove all networks
echo "Removing all networks..."
docker network prune --all --force

# Remove all volumes
echo "Removing all volumes..."
docker volume prune --all --force

echo "Remove others..."
docker system prune --all --volumes

echo "Cleanup completed."
