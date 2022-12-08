# apt-offline

```
$ apt-offline download -h
Download

Usage:
  apt-offline download [flags]

Flags:
      --arch string           Architecture (REQUIRED).
      --docker-image string   Docker image (REQUIRED).
      --download-dir string   Download directory (REQUIRED).
      --dpkg-status string    Path to /var/lib/dpkg/status file.
  -h, --help                  help for download

```

```
$ apt-offline download \
    --download-dir build/download \
    --docker-image debian:stretch \
    --arch armhf \
    --dpkg-status testdata/var/lib/dpkg/status alsa-utils
```
