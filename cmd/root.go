package cmd

import (
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/ngyewch/apt-offline/common"
	"github.com/ngyewch/go-clibase"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", appName),
		Short: "apt-offline",
		RunE:  help,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func help(cmd *cobra.Command, args []string) error {
	err := cmd.Help()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	clibase.AddVersionCmd(rootCmd, func() *goVersion.Info {
		return common.VersionInfo
	})
}

func initConfig() {
	slog.SetLevel(slog.InfoLevel)
}
