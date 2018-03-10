package humanize

import (
	"fmt"
	"go/ast"
)

// MapType is the map in the go code
type MapType struct {
	pkg   *Package
	Key   Type
	Value Type
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.Key.String(), m.Value.String())
}

// Package return the map package
func (m *MapType) Package() *Package {
	return m.pkg
}

// Equal check if the type is equal?
func (m *MapType) Equal(t Type) bool {
	v, ok := t.(*MapType)
	if !ok {
		return false
	}
	if !m.pkg.Equal(v.pkg) {
		return false
	}
	return m.Key.Equal(v.Key) && m.Value.Equal(m.Value)
}

func (m *MapType) lateBind() error {
	return lateBind(m.Key, m.Value)
}

func getMap(p *Package, f *File, t *ast.MapType) Type {
	return &MapType{
		pkg:   p,
		Key:   newType(p, f, t.Key),
		Value: newType(p, f, t.Value),
	}
}
