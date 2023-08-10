package downloader

import (
	"archive/tar"
	"bytes"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/ngyewch/apt-offline/resources"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
)

func (d *Downloader) initDocker() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	d.client = client

	err = d.buildImage()
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) getImageName() string {
	return fmt.Sprintf("apt-offline/%s", d.versionCodename)
}

func (d *Downloader) buildImage() error {
	inputBuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)

	subFs, err := fs.Sub(resources.DockerBuildContextFS, "dockerBuildContext")
	if err != nil {
		return err
	}

	vars := &Vars{
		VersionCodename: d.versionCodename,
		Archived:        d.archived,
	}
	err = createTar(tr, subFs, vars)
	if err != nil {
		return err
	}

	err = tr.Close()
	if err != nil {
		return err
	}

	err = d.client.BuildImage(docker.BuildImageOptions{
		Name:         d.getImageName(),
		InputStream:  inputBuf,
		OutputStream: os.Stdout,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) Download(downloadDir string, arch string, packageNames []string) error {
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	absDownloadDir, err := filepath.Abs(downloadDir)
	if err != nil {
		return err
	}

	u, err := user.Current()
	if err != nil {
		return err
	}

	container, err := d.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        d.getImageName(),
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          packageNames,
			Env: []string{
				fmt.Sprintf("ARCH=%s", arch),
				fmt.Sprintf("UID=%s", u.Uid),
				fmt.Sprintf("GID=%s", u.Gid),
			},
		},
		HostConfig: &docker.HostConfig{
			AutoRemove: true,
			Mounts: []docker.HostMount{
				{
					Target: "/workspace/packages",
					Source: absDownloadDir,
					Type:   "bind",
				},
			},
		},
	})
	if err != nil {
		return err
	}

	err = d.client.StartContainer(container.ID, &docker.HostConfig{
		AutoRemove: true,
	})
	if err != nil {
		return err
	}

	err = d.client.Logs(docker.LogsOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return err
	}

	_, err = d.client.WaitContainer(container.ID)
	if err != nil {
		return err
	}

	return nil
}
