#!/bin/bash
set -xae

source server.env

podman compose -f podman-compose-local.yaml $*
