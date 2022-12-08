package downloader

import docker "github.com/fsouza/go-dockerclient"

type Downloader struct {
	client *docker.Client
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Init() error {
	err := d.initDocker()
	if err != nil {
		return err
	}
	return nil
}
