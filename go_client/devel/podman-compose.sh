#!/bin/bash
set -xae

source go-client.env

podman compose -f podman-compose-local.yaml $*
