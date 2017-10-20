package humanize

import (
	"fmt"
	"go/ast"
	"strconv"
)

// ArrayType is the base array
type ArrayType struct {
	pkg   *Package
	Slice bool
	Len   int
	Type  Type
}

// EllipsisType is slice type but with ...type definition
type EllipsisType struct {
	*ArrayType
}

func (a *ArrayType) Equal(t Type) bool {
	v, ok := t.(*ArrayType)
	if !ok {
		return false
	}

	if !a.pkg.Equal(v.pkg) {
		return false
	}

	if a.Slice != v.Slice {
		return false
	}

	if a.Len != v.Len {
		return false
	}

	return a.Type.Equal(v.Type)
}

func (a *ArrayType) String() string {
	if a.Slice {
		return "[]" + a.Type.String()
	}
	return fmt.Sprintf("[%d]%s", a.Len, a.Type.String())
}

func (a *ArrayType) Package() *Package {
	return a.Package()
}

func (a *ArrayType) lateBind() error {
	return lateBind(a.Type)
}

func (e *EllipsisType) String() string {
	return fmt.Sprintf("[...]%s{}", e.Type.String())
}

func (e EllipsisType) Equal(t Type) bool {
	if v, ok := t.(*EllipsisType); ok {
		return e.ArrayType.Equal(v.ArrayType)
	}
	return false
}

func getArray(p *Package, f *File, t *ast.ArrayType) Type {
	slice := t.Len == nil
	ellipsis := false
	l := 0
	if !slice {
		var (
			ls string
		)
		switch t.Len.(type) {
		case *ast.BasicLit:
			ls = t.Len.(*ast.BasicLit).Value
		case *ast.Ellipsis:
			ls = "0"
			ellipsis = true
		}
		l, _ = strconv.Atoi(ls)
	}
	var at Type
	at = &ArrayType{
		pkg:   p,
		Slice: t.Len == nil,
		Len:   l,
		Type:  newType(p, f, t.Elt),
	}
	if ellipsis {
		at = &EllipsisType{ArrayType: at.(*ArrayType)}
	}
	return at
}
