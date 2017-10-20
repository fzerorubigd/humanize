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
