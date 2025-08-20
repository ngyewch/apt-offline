package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ngyewch/apt-offline/downloader"
	"github.com/ngyewch/apt-offline/dpkg"
	"github.com/urfave/cli/v3"
)

var (
	version string

	flagDownloadDir = &cli.StringFlag{
		Name:     "download-dir",
		Usage:    "download directory",
		Required: true,
	}
	flagVersionCodename = &cli.StringFlag{
		Name:     "version-codename",
		Usage:    "Debian version codename",
		Required: true,
	}
	flagArch = &cli.StringFlag{
		Name:     "arch",
		Usage:    "architecture",
		Required: true,
	}
	flagDpkgStatus = &cli.StringFlag{
		Name:  "dpkg-status",
		Usage: "Path to /var/lib/dpkg/status file",
	}
	flagArchived = &cli.BoolFlag{
		Name:  "archived",
		Usage: "archived mode",
	}

	app = &cli.Command{
		Name:    "apt-offline",
		Usage:   "apt-offline",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "download",
				Usage:  "download",
				Action: doDownload,
				Flags: []cli.Flag{
					flagDownloadDir,
					flagVersionCodename,
					flagArch,
					flagDpkgStatus,
					flagArchived,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doDownload(ctx context.Context, cmd *cli.Command) error {
	downloadDir := cmd.String(flagDownloadDir.Name)
	versionCodename := cmd.String(flagVersionCodename.Name)
	arch := cmd.String(flagArch.Name)
	dpkgStatus := cmd.String(flagDpkgStatus.Name)
	archived := cmd.Bool(flagArchived.Name)

	d := downloader.NewDownloader(versionCodename, archived)

	err := d.Init()
	if err != nil {
		return err
	}

	err = d.Download(downloadDir, arch, cmd.Args().Slice())
	if err != nil {
		return err
	}

	if dpkgStatus != "" {
		f, err := os.Open(dpkgStatus)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		packageStatuses, err := dpkg.ParsePackageStatuses(f)
		if err != nil {
			return err
		}

		dirEntries, err := os.ReadDir(downloadDir)
		if err != nil {
			return err
		}

		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				continue
			}
			if !strings.HasSuffix(dirEntry.Name(), ".deb") {
				continue
			}
			parts := strings.Split(dirEntry.Name()[0:len(dirEntry.Name())-4], "_")
			packageStatus := packageStatuses.FindPackageStatus(parts[0])
			if (packageStatus != nil) && (packageStatus.Status == "install ok installed") && (packageStatus.Architecture == arch) {
				err = os.Remove(filepath.Join(downloadDir, dirEntry.Name()))
				if err != nil {
					return err
				}
				fmt.Printf("# removed %s\n", dirEntry.Name())
			}
		}
	}

	return nil
}
