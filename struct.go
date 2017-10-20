package humanize

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// Field is a single field of a structure, a variable, with tag
type Field struct {
	Name string
	Type Type
	Docs Docs
	Tags reflect.StructTag
}

// Embeds is the embeded type in the struct or interface
type Embed struct {
	Type
	Docs Docs
	Tags reflect.StructTag
}

type Embeds []*Embed

func (e Embeds) Equal(t Embeds) bool {
	for i := range e {
		for j := range t {
			if e[i].Equal(t[j]) {
				continue
			}
		}

		return false
	}
	return true
}

type Fields []*Field

func (f Fields) Equal(t Fields) bool {
	for i := range f {
		for j := range t {
			if f[i].Name == t[j].Name && f[i].Type.Equal(t[j].Type) && f[i].Tags == t[j].Tags {
				continue
			}
		}

		return false
	}
	return true

}

type StructType struct {
	pkg    *Package
	Fields Fields
	Embeds Embeds
}

func (s *StructType) String() string {
	if len(s.Embeds) == 0 && len(s.Fields) == 0 {
		return "struct{}"
	}
	res := "struct {\n"
	for e := range s.Embeds {
		res += "\t" + s.Embeds[e].String() + "\n"
	}

	for f := range s.Fields {
		tags := strings.Trim(string(s.Fields[f].Tags), "`")
		if tags != "" {
			tags = "`" + tags + "`"
		}
		res += fmt.Sprintf("\t%s %s %s\n", s.Fields[f].Name, s.Fields[f].Type.String(), tags)
	}
	return res + "}"
}

func (s *StructType) Package() *Package {
	return s.pkg
}

func (s *StructType) lateBind() error {
	for i := range s.Fields {
		if err := lateBind(s.Fields[i].Type); err != nil {
			return err
		}
	}

	for i := range s.Embeds {
		if err := lateBind(s.Embeds[i].Type); err != nil {
			return err
		}
	}

	return nil
}

func (s *StructType) Equal(t Type) bool {
	v, ok := t.(*StructType)
	if !ok {
		return false
	}

	if !s.pkg.Equal(v.pkg) {
		return false
	}

	if !s.Embeds.Equal(v.Embeds) {
		return false
	}

	if !s.Fields.Equal(v.Fields) {
		return false
	}

	return true
}

func getStruct(p *Package, f *File, t *ast.StructType) Type {
	res := &StructType{
		pkg: p,
	}
	for _, s := range t.Fields.List {
		if s.Names != nil {
			for i := range s.Names {

				f := Field{
					Name: nameFromIdent(s.Names[i]),
					Type: newType(p, f, s.Type),
				}
				if s.Tag != nil {
					f.Tags = reflect.StructTag(s.Tag.Value)
					f.Tags = f.Tags[1 : len(f.Tags)-1]
				}
				f.Docs = docsFromNodeDoc(s.Doc)
				res.Fields = append(res.Fields, &f)
			}
		} else {
			e := Embed{
				Type: newType(p, f, s.Type),
			}
			if s.Tag != nil {
				e.Tags = reflect.StructTag(s.Tag.Value)
				e.Tags = e.Tags[1 : len(e.Tags)-1]
			}
			e.Docs = docsFromNodeDoc(s.Doc)
			res.Embeds = append(res.Embeds, &e)
		}
	}

	return res
}
