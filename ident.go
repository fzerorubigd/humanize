package humanize

import "go/ast"

// IdentType is the normal type name
type IdentType struct {
	pkg   *Package
	Ident string
}

func (i *IdentType) String() string {
	return i.Ident
}

func (i *IdentType) Package() *Package {
	return i.pkg
}

func (i *IdentType) Equal(t Type) bool {
	v, ok := t.(*IdentType)
	if !ok {
		return false
	}
	if !i.pkg.Equal(v.pkg) {
		return false
	}
	return i.Ident == v.Ident
}

func (i *IdentType) lateBind() error {
	if i.pkg == nil {
		i.pkg = getBuiltin()
	}
	return nil
}

func getIdent(p *Package, _ *File, t *ast.Ident) Type {
	// ident is the simplest one.
	return &IdentType{
		pkg:   p,
		Ident: nameFromIdent(t),
	}
}

func getBasicIdent(t string) Type {
	return &IdentType{
		Ident: t,
	}
}

func nameFromIdent(i *ast.Ident) (name string) {
	if i != nil {
		name = i.String()
	}
	return
}
