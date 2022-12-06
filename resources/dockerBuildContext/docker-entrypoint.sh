#!/bin/bash

set -e

apt-get install --download-only -o Dir::Cache="./packages/" -o Dir::Cache::archives="./packages" -y --no-install-recommends "$@"
rm -rf packages/partial packages/lock
