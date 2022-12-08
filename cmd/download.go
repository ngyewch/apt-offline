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
	d := downloader.NewDownloader()

	err := d.Init()
	if err != nil {
		return err
	}

	downloadDir, err := cmd.Flags().GetString("download-dir")
	if err != nil {
		return err
	}

	err = d.Download(downloadDir, args)
	if err != nil {
		return err
	}

	f, err := os.Open("testdata/var/lib/dpkg/status")
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
		if packageStatus != nil {
			fmt.Printf("%s exists\n", packageStatus.Package)
			err = os.Remove(filepath.Join(downloadDir, dirEntry.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().String("download-dir", "", "Download directory.")
	_ = downloadCmd.MarkFlagRequired("download-dir")
}
