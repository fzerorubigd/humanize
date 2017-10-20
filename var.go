package humanize

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Variable is a string represent of a function parameter
type Variable struct {
	pkg  *Package
	Name string
	Type Type
	Docs Docs

	// Any indirect variable need this to determine the type
	file   *File
	caller *ast.CallExpr
	index  int
}

func (v *Variable) String() string {
	return v.Name + " " + v.Type.String()
}

func (v *Variable) Equal(t *Variable) bool {
	if v.pkg.Path != t.pkg.Path {
		return false
	}

	if v.Name != t.Name {
		return false
	}

	return v.Type.Equal(t.Type)
}

func (v *Variable) lateBind() error {
	if v.caller != nil {
		switch c := v.caller.Fun.(type) {
		case *ast.Ident:
			name := nameFromIdent(c)
			bl, err := getBuiltin().FindFunction(name)
			if err == nil {
				v.Type = bl.Type
			} else {
				var t Type
				fn, err := v.pkg.FindFunction(name)
				if err == nil {
					if len(fn.Type.Results) <= v.index {
						return fmt.Errorf("%d result is available but want the %d", len(fn.Type.Results), v.index)
					}
					t = fn.Type.Results[v.index].Type
				} else {
					t, err = checkTypeCast(v.pkg, getBuiltin(), v.caller.Args, name)
					if err != nil {
						return err
					}
				}

				v.Type = t
			}
		case *ast.SelectorExpr:
			var pkg string
			switch c.X.(type) {
			case *ast.Ident:
				pkg = nameFromIdent(c.X.(*ast.Ident))
			case *ast.CallExpr: // TODO : Don't know why, no time for check
				break
			}

			typ := nameFromIdent(c.Sel)
			imprt, err := v.pkg.FindImport(pkg)
			if err != nil {
				// TODO : package currently is not capable of parsing build tags. so ignore this :/
				break
			}
			pkgDef, err := ParsePackage(imprt.Path)
			if err != nil {
				return err
			}
			var t Type
			fn, err := pkgDef.FindFunction(typ)
			if err == nil {
				if len(fn.Type.Results) <= v.index {
					return fmt.Errorf("%d result is available but want %d", len(fn.Type.Results), v.index)
				}
				t = fn.Type.Results[v.index].Type
			} else {
				t, err = checkTypeCast(pkgDef, builtin, v.caller.Args, typ)
				if err != nil {
					return err
				}
			}

			foreignTyp := t
			star := false
			if sType, ok := foreignTyp.(*StarType); ok {
				foreignTyp = sType.Target
				star = true
			}
			switch ft := foreignTyp.(type) {
			case *IdentType:
				// this is a simple hack. if the type is begin with
				// upper case, then its type on that package,
				// else its a global type
				name := ft.Ident
				c := name[0]
				if c >= 'A' && c <= 'Z' {
					if star {
						foreignTyp = &StarType{
							Target: foreignTyp,
						}
					}
					v.Type = &SelectorType{
						selector: pkg,
						file:     v.file,
						Type:     foreignTyp,
					}
				} else {
					if star {
						foreignTyp = &StarType{
							Target: foreignTyp,
						}
					}
					v.Type = foreignTyp
				}

			default:
				// the type is foreign to that package too
				v.Type = ft
			}
		}
	}
	return lateBind(v.Type)
}

func variableFromValue(p *Package, f *File, name string, index int, e []ast.Expr) *Variable {
	var t Type
	var caller *ast.CallExpr
	var ok bool
	first := e[0]
	// if the caller is a CallExpr, then late bind will take care of it
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[index].(type) {
		case *ast.CompositeLit:
			t = newType(p, f, data.Type)
		case *ast.BasicLit:
			switch data.Kind {
			case token.INT:
				t = getBasicIdent("int")
			case token.FLOAT:
				t = getBasicIdent("float64")
			case token.IMAG:
				t = getBasicIdent("complex64")
			case token.CHAR:
				t = getBasicIdent("char")
			case token.STRING:
				t = getBasicIdent("string")
			}
			//default:
			//fmt.Printf("var value => %T", e[index])
			//fmt.Printf("%s", src[data.Pos()-1:data.End()-1])
		}
	}
	return &Variable{
		pkg:    p,
		Name:   name,
		Type:   t,
		caller: caller,
		index:  index,
		file:   f,
	}
}

func variableFromExpr(p *Package, f *File, name string, e ast.Expr) *Variable {
	return &Variable{
		pkg:  p,
		Name: name,
		Type: newType(p, f, e),
	}
}

// newVariable return an array of variables in the scope
func newVariable(p *Package, f *File, v *ast.ValueSpec, c *ast.CommentGroup) []*Variable {
	var res []*Variable
	for i := range v.Names {
		name := nameFromIdent(v.Names[i])
		var n *Variable
		if v.Type != nil {
			n = variableFromExpr(p, f, name, v.Type)
		} else {
			if len(v.Values) != 0 {
				n = variableFromValue(p, f, name, i, v.Values)
			}
		}
		n.Docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}
