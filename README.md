# apt-offline

Offline package management for Debian. It can be used to download packages and their dependencies to be installed later on a disconnected machine.

## Pre-requisites

* Docker

## Installation

```
go install github.com/ngyewch/apt-offline@latest
```

## Usage

```
$ apt-offline download -h
Download

Usage:
  apt-offline download [flags]

Flags:
      --arch string               Architecture (REQUIRED).
      --archived                  Archived mode.
      --download-dir string       Download directory (REQUIRED).
      --dpkg-status string        Path to /var/lib/dpkg/status file.
  -h, --help                      help for download
      --version-codename string   Debian version codename (REQUIRED).
```

```
$ apt-offline download \
    --download-dir build/download \
    --version-codename stretch \
    --arch armhf \
    --dpkg-status testdata/var/lib/dpkg/status \
    --archived \
    alsa-utils rsync sudo
```
