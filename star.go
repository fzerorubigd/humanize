package humanize

import "go/ast"

// StarType is the pointer of a type
type StarType struct {
	Target Type
	pkg    *Package
	file   *File
}

func (s *StarType) String() string {
	return "*" + s.Target.String()
}

// Package get the package name
func (s *StarType) Package() *Package {
	return s.pkg
}

// Equal check if two star type are equal
func (s *StarType) Equal(t Type) bool {
	v, ok := t.(*StarType)
	if !ok {
		return false
	}
	return s.Target.Equal(v.Target)
}

func (s *StarType) lateBind() error {
	return lateBind(s.Target)
}

func getStar(p *Package, f *File, t *ast.StarExpr) Type {
	return &StarType{
		Target: newType(p, f, t.X),
		pkg:    p,
		file:   f,
	}
}
