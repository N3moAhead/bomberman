#!/bin/bash
set -xae

source server.env

COMPOSE_FILE="podman-compose-local.yaml"
CONTAINER_NAME="bomberman-sever"
SERVICE_NAME="bomberman-server"

# Check if the container is already running
if podman ps --filter "name=${CONTAINER_NAME}" --format "{{.Names}}" | grep -q "${CONTAINER_NAME}"; then
  echo "An instance of ${CONTAINER_NAME} is already running. Showing logs..."
  podman compose -f ${COMPOSE_FILE} logs -f ${SERVICE_NAME}
else
  echo "Starting a new instance of ${CONTAINER_NAME}..."
  podman compose -f ${COMPOSE_FILE} $*
fi
