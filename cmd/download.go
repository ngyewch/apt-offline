package cmd

import (
	"fmt"
	"github.com/ngyewch/apt-offline/downloader"
	"github.com/ngyewch/apt-offline/dpkg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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
	dpkgStatus, err := cmd.Flags().GetString("dpkg-status")
	if err != nil {
		return err
	}

	downloadDir, err := cmd.Flags().GetString("download-dir")
	if err != nil {
		return err
	}

	versionCodename, err := cmd.Flags().GetString("version-codename")
	if err != nil {
		return err
	}

	arch, err := cmd.Flags().GetString("arch")
	if err != nil {
		return err
	}

	archived, err := cmd.Flags().GetBool("archived")
	if err != nil {
		return err
	}

	d := downloader.NewDownloader(versionCodename, archived)

	err = d.Init()
	if err != nil {
		return err
	}

	err = d.Download(downloadDir, arch, args)
	if err != nil {
		return err
	}

	if dpkgStatus != "" {
		f, err := os.Open(dpkgStatus)
		if err != nil {
			return err
		}
		defer f.Close()

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

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().String("download-dir", "", "Download directory (REQUIRED).")
	downloadCmd.Flags().String("version-codename", "", "Debian version codename (REQUIRED).")
	downloadCmd.Flags().String("arch", "", "Architecture (REQUIRED).")
	downloadCmd.Flags().String("dpkg-status", "", "Path to /var/lib/dpkg/status file.")
	downloadCmd.Flags().Bool("archived", false, "Archived mode.")

	_ = downloadCmd.MarkFlagRequired("download-dir")
	_ = downloadCmd.MarkFlagRequired("version-codename")
	_ = downloadCmd.MarkFlagRequired("arch")
}
