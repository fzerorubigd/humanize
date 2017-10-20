package humanize

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Import is one imported path
type Import struct {
	Package   string
	Canonical string
	Path      string
	Docs      Docs

	Folder string

	pkg *Package
}

func (i *Import) String() string {
	if i.Canonical != i.Package {
		return fmt.Sprintf(`%s "%s"`, i.Canonical, i.Path)
	}

	return fmt.Sprintf(`"%s"`, i.Path)
}

type importWalker struct {
	pkgName  string
	resolved string
}

func (iw *importWalker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.File:
			iw.pkgName = nameFromIdent(t.Name)
		default:
		}
	}

	return iw
}

// LoadPackage is the function to load import package
func (i Import) LoadPackage() (*Package, error) {
	if i.Folder == "" {
		return nil, fmt.Errorf("the package '%s' is not resolved", i.Path)
	}
	if i.pkg == nil {
		p, err := parsePackageFullPath(i.Path, i.Folder)
		if err != nil {
			return nil, err
		}
		i.pkg = p
	}
	return i.pkg, nil

}

func peekPackageName(pkg string, base ...string) (string, string) {
	_, name := filepath.Split(pkg)
	folder, err := translateToFullPath(pkg, base...)
	if err != nil {
		return name, ""
	}
	iw := &importWalker{}
	err = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, err error) error {
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fset := token.NewFileSet()
			fle, err := parser.ParseFile(fset, "", data, parser.PackageClauseOnly)
			if err != nil {
				return nil // try another file?
			}
			iw.resolved = folder
			ast.Walk(iw, fle)
			// no need to continue
			return filepath.SkipDir
		},
	)
	resolved := ""
	if iw.pkgName != "" {
		name = iw.pkgName
		resolved = iw.resolved

	}
	// can not parse it, use the folder name
	return name, resolved
}

// newImport extract a new import entry
func newImport(p *Package, f *File, i *ast.ImportSpec, c *ast.CommentGroup) *Import {
	res := &Import{
		Package: "",
		Path:    strings.Trim(i.Path.Value, `"`),
		Docs:    docsFromNodeDoc(c, i.Doc),
	}
	if i.Name != nil {
		res.Canonical = i.Name.String()
	}
	res.Package, res.Folder = peekPackageName(res.Path, p.Path)
	if res.Canonical == "" {
		res.Canonical = res.Package
	}

	return res
}
