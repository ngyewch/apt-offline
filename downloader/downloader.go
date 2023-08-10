package downloader

import docker "github.com/fsouza/go-dockerclient"

type Downloader struct {
	versionCodename string
	archived        bool
	client          *docker.Client
}

type Vars struct {
	VersionCodename string
	Archived        bool
}

func NewDownloader(versionCodename string, archived bool) *Downloader {
	return &Downloader{
		versionCodename: versionCodename,
		archived:        archived,
	}
}

func (d *Downloader) Init() error {
	err := d.initDocker()
	if err != nil {
		return err
	}
	return nil
}
