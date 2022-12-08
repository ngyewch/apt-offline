package downloader

import docker "github.com/fsouza/go-dockerclient"

type Downloader struct {
	image  string
	client *docker.Client
}

type Vars struct {
	Image string
}

func NewDownloader(image string) *Downloader {
	return &Downloader{
		image: image,
	}
}

func (d *Downloader) Init() error {
	err := d.initDocker()
	if err != nil {
		return err
	}
	return nil
}
