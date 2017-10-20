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

func (m *MapType) Package() *Package {
	return m.pkg
}

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
	if err := lateBind(m.Key); err != nil {
		return err
	}

	return lateBind(m.Value)
}

func getMap(p *Package, f *File, t *ast.MapType) Type {
	return &MapType{
		pkg:   p,
		Key:   newType(p, f, t.Key),
		Value: newType(p, f, t.Value),
	}
}
