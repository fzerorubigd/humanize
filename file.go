package humanize

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// File is a single file in a structure
type File struct {
	FileName    string
	PackageName string

	Docs      Docs
	Imports   []*Import
	Variables []*Variable
	Functions []*Function
	Constants []*Constant
	Types     []*TypeName
}

func (f *File) lateBind() error {
	for i := range f.Variables {
		if err := f.Variables[i].lateBind(); err != nil {
			return err
		}
	}

	for i := range f.Functions {
		if err := f.Functions[i].lateBind(); err != nil {
			return err
		}
	}

	for i := range f.Constants {
		if err := f.Constants[i].lateBind(); err != nil {
			return err
		}
	}

	for i := range f.Types {
		if err := f.Types[i].lateBind(); err != nil {
			return err
		}
	}

	return nil
}

type walker struct {
	src     string
	File    *File
	Package *Package
}

func (fv *walker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.File:
			fv.File.PackageName = nameFromIdent(t.Name)
			fv.File.Docs = docsFromNodeDoc(t.Doc)
		case *ast.FuncDecl:
			fv.File.Functions = append(fv.File.Functions, getFunction(fv.Package, fv.File, t))
			return nil // Do not go deeper
		case *ast.GenDecl:
			for i := range t.Specs {
				switch decl := t.Specs[i].(type) {
				case *ast.ImportSpec:
					fv.File.Imports = append(fv.File.Imports, newImport(fv.Package, fv.File, decl, t.Doc))
				case *ast.ValueSpec:
					if t.Tok.String() == "var" {
						fv.File.Variables = append(fv.File.Variables, newVariable(fv.Package, fv.File, decl, t.Doc)...)
					} else if t.Tok.String() == "const" {
						var last *Constant
						if len(fv.File.Constants) > 0 {
							last = fv.File.Constants[len(fv.File.Constants)-1]
						}
						fv.File.Constants = append(fv.File.Constants, newConstant(fv.Package, fv.File, decl, t.Doc, last)...)
					}
				case *ast.TypeSpec:
					fv.File.Types = append(fv.File.Types, newTypeName(fv.Package, fv.File, decl, t.Doc))
				}
			}
			return nil
		default:
			//fmt.Printf("\n%T=====>%+v", t, t)
		}
	}
	return fv
}

// ParseFile try to parse a single file for its annotations
func ParseFile(src string, p *Package) (*File, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fw := &walker{}
	fw.src = src
	fw.File = &File{}
	fw.Package = p

	ast.Walk(fw, f)

	return fw.File, nil
}
