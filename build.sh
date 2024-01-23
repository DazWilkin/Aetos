#!/usr/bin/env bash

CHECKSUM=$(git rev-parse HEAD)
VERSION="0.0.1"

IMAGE="localhost:32000/aetos"

podman build \
--tag="${IMAGE}:${CHECKSUM}" \
--build-arg=CHECKSUM="${CHECKSUM}" \
--build-arg=VERSION="${VERSION}" \
--file=./Dockerfile \
${PWD}

podman push ${IMAGE}:${CHECKSUM}