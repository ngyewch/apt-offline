#!/bin/bash

set -e

VALID_ARGS=$(getopt -o a --long archived -- "$@")
if [[ $? -ne 0 ]]; then
    exit 1;
fi

ARCHIVED=0
eval set -- "$VALID_ARGS"
while [ : ]; do
  case "$1" in
    -a | --archived)
        ARCHIVED=1
        shift
        ;;
    --) shift;
        break
        ;;
  esac
done

if [ -f /etc/os-release ]; then
  source /etc/os-release
  if [ "${ID}" == "debian" ]; then
    if [ "${ARCHIVED}" == "1" ]; then
      cat <<EOT > /etc/apt/sources.list
deb [trusted=yes] http://archive.debian.org/debian ${VERSION_CODENAME} main non-free contrib
deb-src [trusted=yes] http://archive.debian.org/debian/ ${VERSION_CODENAME} main non-free contrib
deb [trusted=yes] http://archive.debian.org/debian-security/ ${VERSION_CODENAME}/updates main non-free contrib
EOT
    fi
  fi
fi

dpkg --add-architecture ${ARCH}
apt update -o APT::Architecture="${ARCH}" -o APT::Architectures="${ARCH}"
apt-get install --download-only -o Dir::Cache="./packages/" -o Dir::Cache::archives="./packages" -y --no-install-recommends "$@"
rm -rf packages/partial packages/lock
chown -R ${UID}:${GID} packages/*
