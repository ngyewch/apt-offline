package cmd

import (
	"archive/tar"
	"bytes"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/ngyewch/apt-offline/resources"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"time"
)

var (
	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download",
		Args:  cobra.MinimumNArgs(1),
		RunE:  download,
	}
)

func download(cmd *cobra.Command, args []string) error {
	err := buildImage(args)
	if err != nil {
		return err
	}
	return nil
}

func buildImage(packageNames []string) error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	inputBuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)

	subFs, err := fs.Sub(resources.DockerBuildContextFS, "dockerBuildContext")
	if err != nil {
		return err
	}

	err = createTar(tr, subFs)
	if err != nil {
		return err
	}

	err = tr.Close()
	if err != nil {
		return err
	}

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         "apt-offline:latest",
		InputStream:  inputBuf,
		OutputStream: os.Stdout,
	})
	if err != nil {
		return err
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        "apt-offline:latest",
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          packageNames,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove: true,
		},
	})
	if err != nil {
		return err
	}

	err = client.StartContainer(container.ID, &docker.HostConfig{
		AutoRemove: true,
	})
	if err != nil {
		return err
	}

	err = client.Logs(docker.LogsOptions{
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

	_, err = client.WaitContainer(container.ID)
	if err != nil {
		return err
	}

	return nil
}

func createTar(tr *tar.Writer, filesystem fs.FS) error {
	return fs.WalkDir(filesystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		fi, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			t := time.Now()
			err = tr.WriteHeader(&tar.Header{
				Name:       path,
				Size:       0,
				Mode:       int64(fi.Mode()),
				ModTime:    t,
				AccessTime: t,
				ChangeTime: t,
			})
			if err != nil {
				return err
			}
		} else {
			f, err := filesystem.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			contentBytes, err := io.ReadAll(f)
			t := time.Now()
			err = tr.WriteHeader(&tar.Header{
				Name:       path,
				Size:       int64(len(contentBytes)),
				Mode:       int64(fi.Mode()),
				ModTime:    t,
				AccessTime: t,
				ChangeTime: t,
			})
			if err != nil {
				return err
			}
			_, err = tr.Write(contentBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
