#!/bin/bash

set -e

dpkg --add-architecture ${ARCH}
apt update -o APT::Architecture="${ARCH}" -o APT::Architectures="${ARCH}"
apt-get install --download-only -o Dir::Cache="./packages/" -o Dir::Cache::archives="./packages" -y --no-install-recommends "$@"
rm -rf packages/partial packages/lock
chown -R ${UID}:${GID} packages/*
