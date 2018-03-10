package humanize

import (
	"fmt"
	"go/ast"
)

// SelectorType is the type in another package
type SelectorType struct {
	Type Type

	file     *File
	selector string
	typeName *TypeName
	imp      *Import
}

func (s *SelectorType) String() string {
	return s.selector + "." + s.Type.String()
}

// Package is the package of selector
func (s *SelectorType) Package() *Package {
	return s.Type.Package()
}

// Equal if two types are equal
func (s *SelectorType) Equal(t Type) bool {
	if b := s.Type.Equal(t); b {
		return b
	}

	v, ok := t.(*SelectorType)
	if !ok {
		return false
	}

	return s.Type.Equal(v.Type)
}

// TargetPackage is the selector target package
func (s *SelectorType) TargetPackage() (*Package, error) {
	if s.imp != nil {
		return s.imp.LoadPackage()
	}

	return nil, fmt.Errorf("the package is not loadable")
}

func (s *SelectorType) lateBind() error {
	var err error
	s.imp, err = s.file.FindImport(s.selector)
	if err != nil {
		return err
	}

	p2, err := s.imp.LoadPackage()
	if err != nil {
		return err
	}
	switch t := s.Type.(type) {
	case *IdentType:
		t.pkg = p2
		s.typeName, err = p2.FindType(t.Ident)
	case *StarType:
		t.Target.(*IdentType).pkg = p2
		s.typeName, err = p2.FindType(t.Target.(*IdentType).Ident)
	}
	if err != nil {
		return err
	}
	return s.typeName.lateBind()
}

func getSelector(p *Package, f *File, t *ast.SelectorExpr) Type {
	switch it := t.X.(type) {
	case *ast.Ident:
		res := &SelectorType{
			Type:     getIdent(p, f, t.Sel).(*IdentType),
			selector: nameFromIdent(it),
			file:     f,
		}
		return res
	default:
		panic(fmt.Sprintf("%T is not supported. please report this (with sample code) to add support for it", it))
	}
}
