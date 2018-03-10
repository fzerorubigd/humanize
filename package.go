package humanize

import (
	"os"
	"path/filepath"
)

// Package is the package in one place
type Package struct {
	Path string
	Name string

	Files []*File
}

// Bind try to bind package. this check every type and bind them to actual value
func (p *Package) Bind() error {
	for i := range p.Files {
		if err := p.Files[i].lateBind(); err != nil {
			return err
		}
	}

	return nil
}

// Equal check if to packages are equal
func (p *Package) Equal(t *Package) bool {
	return p.Path == t.Path
}

func parsePackageFullPath(path, folder string) (*Package, error) {
	if p := getCache(folder); p != nil {
		return p, nil
	}

	var (
		p = &Package{}
		e error
	)
	p.Path = path
	e = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, _ error) error {
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fl, err := ParseFile(data, p)
			if err != nil {
				return err
			}
			fl.FileName = path
			p.Files = append(p.Files, fl)
			if p.Name == "" {
				p.Name = fl.PackageName
			}

			return nil
		},
	)
	if e != nil {
		return nil, e
	}
	setCache(folder, p)

	return p, nil
}

// ParsePackage is here for loading a single package and parse all files in it
// if the package is imported from another package, the other parameter is required for
// checking vendors of that package.
func ParsePackage(path string, packages ...string) (*Package, error) {
	folder, err := translateToFullPath(path, packages...)
	if err != nil {
		return nil, err
	}

	return parsePackageFullPath(path, folder)
}
