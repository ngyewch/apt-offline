package dpkg

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type PackageStatuses struct {
	Entries []*PackageStatus
}

type PackageStatus struct {
	Package       string
	Essential     bool
	Status        string
	Priority      string
	Section       string
	InstalledSize int64
	Origin        string
	Maintainer    string
	Bugs          string
	Architecture  string
	MultiArch     string
	Source        string
	Version       string
	Replaces      []*PackageDep
	Provides      []*PackageDep
	Depends       []*PackageDep
	PreDepends    []*PackageDep
	Recommends    []*PackageDep
	Suggests      []*PackageDep
	Breaks        []*PackageDep
	Enhances      []*PackageDep
	Conflicts     []*PackageDep
	ConfFiles     []*ConfFile
	Description   string
	Homepage      string
}

type PackageDep struct {
	Name       string
	Constraint *semver.Constraints
}

type ConfFile struct {
	Path     string
	Checksum string
}

func ParsePackageStatuses(r io.Reader) (*PackageStatuses, error) {
	packageStatusList := make([]*PackageStatus, 0)
	var packageStatus *PackageStatus = &PackageStatus{}
	err := parseStatusFile(r, func(lineNo int, key string, value string) error {
		switch key {
		case "Package":
			packageStatus.Package = value
			break
		case "Essential":
			if value == "yes" {
				packageStatus.Essential = true
			} else if value == "no" {
				packageStatus.Essential = false
			} else {
				return fmt.Errorf("error at line %d", lineNo)
			}
			break
		case "Status":
			packageStatus.Status = value
			break
		case "Priority":
			packageStatus.Priority = value
			break
		case "Section":
			packageStatus.Section = value
			break
		case "Installed-Size":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error at line %d", lineNo)
			}
			packageStatus.InstalledSize = v
			break
		case "Origin":
			packageStatus.Origin = value
			break
		case "Maintainer":
			packageStatus.Maintainer = value
			break
		case "Bugs":
			packageStatus.Bugs = value
			break
		case "Architecture":
			packageStatus.Architecture = value
			break
		case "Multi-Arch":
			packageStatus.MultiArch = value
			break
		case "Source":
			packageStatus.Source = value
			break
		case "Version":
			packageStatus.Version = value
			break
		case "Replaces":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Replaces = packageDeps
			break
		case "Provides":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Provides = packageDeps
			break
		case "Depends":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Depends = packageDeps
			break
		case "Pre-Depends":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.PreDepends = packageDeps
			break
		case "Recommends":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Recommends = packageDeps
			break
		case "Suggests":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Suggests = packageDeps
			break
		case "Breaks":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Breaks = packageDeps
			break
		case "Enhances":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Enhances = packageDeps
			break
		case "Conflicts":
			packageDeps, err := parsePackageDeps(value)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", lineNo, err)
			}
			packageStatus.Conflicts = packageDeps
			break
		case "Conffiles":
			confFiles := parseConfFiles(value)
			packageStatus.ConfFiles = confFiles
			break
		case "Description":
			packageStatus.Description = value
			break
		case "Homepage":
			packageStatus.Homepage = value
			break
		default:
			return fmt.Errorf("unknown key '%s' at line %d", key, lineNo)
		}
		return nil
	}, func(lineNo int) error {
		packageStatusList = append(packageStatusList, packageStatus)
		packageStatus = &PackageStatus{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &PackageStatuses{
		Entries: packageStatusList,
	}, nil
}

func (p *PackageStatuses) FindPackageStatus(name string) *PackageStatus {
	for _, packageStatus := range p.Entries {
		if packageStatus.Package == name {
			return packageStatus
		}
	}
	return nil
}

func parsePackageDeps(s string) ([]*PackageDep, error) {
	parts := strings.Split(s, ", ")
	packageDeps := make([]*PackageDep, 0)
	for _, part := range parts {
		p := strings.Index(part, " (")
		if p >= 0 {
			constraint, err := semver.NewConstraint(part[p+2 : len(part)-1])
			if err != nil {
				//return nil, err
				// TODO Hack
				packageDeps = append(packageDeps, &PackageDep{
					Name: part,
				})
			} else {
				packageDeps = append(packageDeps, &PackageDep{
					Name:       part,
					Constraint: constraint,
				})
			}
		} else {
			packageDeps = append(packageDeps, &PackageDep{
				Name: part,
			})
		}
	}
	return packageDeps, nil
}

func parseConfFiles(s string) []*ConfFile {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	confFiles := make([]*ConfFile, 0)
	for _, part := range parts {
		subParts := strings.SplitN(part, " ", 2)
		confFiles = append(confFiles, &ConfFile{
			Path:     subParts[0],
			Checksum: subParts[1],
		})
	}
	return confFiles
}
