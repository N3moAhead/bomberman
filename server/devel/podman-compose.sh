#!/bin/bash
set -xae

source server.env

podman compose -p bomberman-server -f podman-compose-local.yaml $*
