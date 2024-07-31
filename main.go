package main

import (
	"fmt"
	"github.com/ngyewch/apt-offline/downloader"
	"github.com/ngyewch/apt-offline/dpkg"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	version         string
	commit          string
	commitTimestamp string

	flagDownloadDir = &cli.PathFlag{
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
	flagDpkgStatus = &cli.PathFlag{
		Name:  "dpkg-status",
		Usage: "Path to /var/lib/dpkg/status file",
	}
	flagArchived = &cli.BoolFlag{
		Name:  "archived",
		Usage: "archived mode",
	}

	app = &cli.App{
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
	cli.VersionPrinter = func(cCtx *cli.Context) {
		var parts []string
		if version != "" {
			parts = append(parts, fmt.Sprintf("version=%s", version))
		}
		if commit != "" {
			parts = append(parts, fmt.Sprintf("commit=%s", commit))
		}
		if commitTimestamp != "" {
			formattedCommitTimestamp := func(commitTimestamp string) string {
				epochSeconds, err := strconv.ParseInt(commitTimestamp, 10, 64)
				if err != nil {
					return ""
				}
				t := time.Unix(epochSeconds, 0)
				return t.Format(time.RFC3339)
			}(commitTimestamp)
			if formattedCommitTimestamp != "" {
				parts = append(parts, fmt.Sprintf("commitTimestamp=%s", formattedCommitTimestamp))
			}
		}
		fmt.Println(strings.Join(parts, " "))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doDownload(cCtx *cli.Context) error {
	downloadDir := flagDownloadDir.Get(cCtx)
	versionCodename := flagVersionCodename.Get(cCtx)
	arch := flagArch.Get(cCtx)
	dpkgStatus := flagDpkgStatus.Get(cCtx)
	archived := flagArchived.Get(cCtx)

	d := downloader.NewDownloader(versionCodename, archived)

	err := d.Init()
	if err != nil {
		return err
	}

	err = d.Download(downloadDir, arch, cCtx.Args().Slice())
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
