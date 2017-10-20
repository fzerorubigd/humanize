package humanize

import (
	"fmt"
	"go/ast"
)

// Type is the interface for all types without name
type Type interface {
	fmt.Stringer
	// Package return the package name of the type
	Package() *Package

	Equal(Type) bool
}

func newType(p *Package, f *File, e ast.Expr) Type {
	switch t := e.(type) {
	case *ast.Ident:
		return getIdent(p, f, t)
	case *ast.StarExpr:
		return getStar(p, f, t)
	case *ast.ArrayType:
		return getArray(p, f, t)
	case *ast.MapType:
		return getMap(p, f, t)
	case *ast.StructType:
		return getStruct(p, f, t)
	case *ast.SelectorExpr:
		return getSelector(p, f, t)
	case *ast.ChanType:
		return getChannel(p, f, t)
	case *ast.FuncType:
		return getFunc(p, f, t)
	case *ast.InterfaceType:
		return getInterface(p, f, t)
	default:
		return nil
	}
}

// FindTypeName try to find type name based on the type. it can fail for
// anonymous types, so watch about the result
func FindTypeName(t Type) (*TypeName, error) {
	switch v := t.(type) {
	case *IdentType:
		return t.Package().FindType(v.Ident)
	case *StarType:
		return FindTypeName(v.Target)
	default:
		return nil, fmt.Errorf("%T is not supported yet", t)
	}
}
