package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	testCmd = &cobra.Command{
		Use:   "test",
		Short: "Test",
		Args:  cobra.ExactArgs(0),
		RunE:  test,
	}
)

func test(cmd *cobra.Command, args []string) error {
	f, err := os.Open("testdata/var/lib/dpkg/status")
	if err != nil {
		return err
	}
	defer f.Close()

	err = parseStatusFile(f, func(key string, value string) error {
		fmt.Printf("%s: %s\n", key, value)
		return nil
	}, func() error {
		fmt.Println("-------------------")
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type PackageStatus struct {
	Package       string
	Status        string
	Priority      string
	Architecture  string
	MultiArch     string
	Maintainer    string
	Version       string
	Section       string
	InstalledSize int64
	Replaces      []string // should be a struct
	Depends       []string // should be a struct
	Suggests      []string // should be a struct
	Breaks        []string // should be a struct
	Conflicts     []string // should be a struct
	Description   string
	Source        string
	Homepage      string
}

func parsePackageStatusFile(r io.Reader) ([]*PackageStatus, error) {
	packageStatusList := make([]*PackageStatus, 0)
	return packageStatusList, nil
}

func parseStatusFile(r io.Reader, entryFunc func(key string, value string) error, commitFunc func() error) error {
	scanner := bufio.NewScanner(r)
	key := ""
	value := ""
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		if line == "" {
			if key != "" {
				err := entryFunc(key, value)
				if err != nil {
					return err
				}
				key = ""
				value = ""
			}
			err := commitFunc()
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(line, " ") {
			if key == "" {
				return fmt.Errorf("invalid continuation at line %d", lineNo)
			}
			value += "\n"
			v := line[1:]
			if v != "." {
				value += v
			}
		} else {
			if key != "" {
				err := entryFunc(key, value)
				if err != nil {
					return err
				}
				key = ""
				value = ""
			}
			p := strings.Index(line, ":")
			key = line[0:p]
			if p+2 < len(line) {
				value = line[p+2:]
			}
		}
	}
	if key != "" {
		err := entryFunc(key, value)
		if err != nil {
			return err
		}
		key = ""
		value = ""
	}
	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
