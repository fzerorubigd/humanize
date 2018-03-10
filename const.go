package humanize

// TODO : iota support is very limited and bad

import (
	"fmt"
	"go/ast"
	"go/token"
)

var (
	lastConst Type
)

// Constant is a string represent of a function parameter
type Constant struct {
	pkg   *Package
	Name  string
	Type  Type
	Docs  Docs
	Value string

	caller *ast.CallExpr
	index  int
}

func (c *Constant) lateBind() error {
	return lateBind(c.Type)
}

// String is the constant in go source code
func (c *Constant) String() string {
	return fmt.Sprintf("%s %s = %s", c.Name, c.Type.String(), c.Value)
}

// Equal check for equality
// TODO : currently it check if the value is equal and the type is equal too.
func (c *Constant) Equal(t *Constant) bool {
	if c.Value != t.Value {
		return false
	}

	return c.Type.Equal(t.Type)
}

func constantFromValue(p *Package, f *File, name string, indx int, e []ast.Expr) *Constant {
	var t Type
	var caller *ast.CallExpr
	var ok bool
	if len(e) == 0 {
		return &Constant{
			pkg:  p,
			Name: name,
		}
	}
	first := e[0]
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[indx].(type) {
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
				//default:
				//fmt.Printf("var value => %T", e[index])
				//fmt.Printf("%s", src[data.Pos()-1:data.End()-1])
			}
		case *ast.Ident:
			t = getIdent(p, f, data)
		}
	}
	return &Constant{
		pkg:    p,
		Name:   name,
		Type:   t,
		caller: caller,
		index:  indx,
	}
}
func constantFromExpr(p *Package, f *File, name string, e ast.Expr) *Constant {
	return &Constant{
		pkg:  p,
		Name: name,
		Type: newType(p, f, e),
	}
}

func getConstantValue(a []ast.Expr, lastVal string) string {
	if len(a) == 0 {
		return lastVal
	}
	switch first := a[0].(type) {
	case *ast.BasicLit:
		return first.Value
	default:
		//fmt.Printf("%T ==> %+v", first, first)
		return "NotSupportedYet"
	}
}

// newConstant return an array of constant in the scope
func newConstant(p *Package, f *File, v *ast.ValueSpec, c *ast.CommentGroup, last *Constant) []*Constant {
	var res []*Constant
	for i := range v.Names {
		name := nameFromIdent(v.Names[i])
		var n *Constant
		if v.Type != nil {
			n = constantFromExpr(p, f, name, v.Type)
		} else {
			n = constantFromValue(p, f, name, i, v.Values)
		}
		l := ""
		if last != nil {
			l = last.Value
		}
		n.Value = getConstantValue(v.Values, l)
		if n.Type == nil {
			n.Type = lastConst
		} else {
			lastConst = n.Type
		}
		n.Name = name
		n.Docs = docsFromNodeDoc(c, v.Doc)
		last = n
		res = append(res, n)
	}

	return res
}
