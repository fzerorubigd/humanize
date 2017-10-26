package humanize

import (
	"go/ast"
)

// TypeName is the type with name in a package
type TypeName struct {
	pkg  *Package
	file *File

	Docs Docs
	Type Type
	Name string

	Methods     []*Function
	StarMethods []*Function
}

func (tn *TypeName) Equal(t *TypeName) bool {
	if !tn.pkg.Equal(t.pkg) {
		return false
	}

	if tn.Name != t.Name {
		return false
	}

	return tn.Type.Equal(t.Type)
}

func (tn *TypeName) InstanceOf(t Type, samePkg bool) (bool, bool) {
	var pointer bool
	if v, ok := t.(*StarType); ok {
		pointer = true
		t = v.Target
	}
	switch tt := t.(type) {
	case *IdentType:
		if samePkg && !tt.pkg.Equal(tn.pkg) {
			return false, pointer
		}
		if tt.Ident != tn.Name {
			return false, pointer
		}
	}

	return true, pointer
}

func (tn *TypeName) String() string {
	return tn.Name + " " + tn.Type.String()
}

func (tn *TypeName) lateBind() error {
	if err := lateBind(tn.Type); err != nil {
		return err
	}
	for i := range tn.pkg.Files {
		for j := range tn.pkg.Files[i].Functions {
			rec := tn.pkg.Files[i].Functions[j].Receiver
			if rec == nil {
				continue
			}

			// method must be in same package with type
			if ok, star := tn.InstanceOf(rec.Type, true); ok {
				if star {
					tn.StarMethods = append(tn.StarMethods, tn.pkg.Files[i].Functions[j])
				} else {
					tn.Methods = append(tn.Methods, tn.pkg.Files[i].Functions[j])
				}
			}
		}
	}

	return nil
}

func (tn *TypeName) Package() *Package {
	return tn.pkg
}

func (tn *TypeName) File() *File {
	return tn.file
}

// newTypeName handle a type with name
func newTypeName(p *Package, f *File, t *ast.TypeSpec, c *ast.CommentGroup) *TypeName {
	doc := docsFromNodeDoc(c, t.Doc)
	return &TypeName{
		pkg:  p,
		file: f,
		Docs: doc,
		Type: newType(p, f, t.Type),
		Name: nameFromIdent(t.Name),
	}
}
