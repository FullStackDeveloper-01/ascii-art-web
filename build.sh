#!/bin/bash
# Build the image, then run a container from it.
# Usage: bash build.sh

set -e

IMAGE_NAME="ascii-art-web"
CONTAINER_NAME="aaw-container"
PORT="8090"

echo "==> Building image: $IMAGE_NAME"
docker image build -f Dockerfile -t "$IMAGE_NAME" .

echo "==> Removing any old container named $CONTAINER_NAME"
docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

echo "==> Running container: $CONTAINER_NAME"
docker container run -p "$PORT:$PORT" --detach --name "$CONTAINER_NAME" "$IMAGE_NAME"

echo "==> Done. Server should be at http://localhost:$PORT"
docker ps -a --filter "name=$CONTAINER_NAME"
