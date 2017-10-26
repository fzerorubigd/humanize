package humanize

import (
	"fmt"
	"go/ast"
)

type Function struct {
	pkg  *Package
	file *File

	Name     string
	Docs     Docs
	Type     *FuncType
	Receiver *Variable // Nil means normal function
}

func (f *Function) String() string {
	s := "func "
	if f.Receiver != nil {
		s += fmt.Sprintf("(%s)", f.Receiver.Type.String())
	}
	s += f.Name + f.Type.Sign()
	return s
}

func (f *Function) Equal(t *Function) bool {
	if !f.pkg.Equal(t.pkg) {
		return false
	}

	if f.Name != t.Name {
		return false
	}

	if f.Receiver == nil && t.Receiver != nil {
		return false
	}
	if f.Receiver != nil && t.Receiver == nil {
		return false
	}

	if f.Receiver != nil {
		if !f.Receiver.Equal(t.Receiver) {
			return false
		}
	}

	return f.Type.Equal(t.Type)
}

func (f *Function) lateBind() error {
	if err := f.Type.lateBind(); err != nil {
		return err
	}

	if f.Receiver != nil {
		return f.Receiver.lateBind()
	}

	return nil
}

func (f *Function) Package() *Package {
	return f.pkg
}

func (f *Function) File() *File {
	return f.file
}

// NewFunction return a single function annotation
func getFunction(p *Package, fl *File, f *ast.FuncDecl) *Function {
	res := &Function{
		pkg:  p,
		file: fl,
		Name: nameFromIdent(f.Name),
		Docs: docsFromNodeDoc(f.Doc),
	}

	if f.Recv != nil {
		// Method receiver is only one parameter
		for i := range f.Recv.List {
			n := ""
			if f.Recv.List[i].Names != nil {
				n = nameFromIdent(f.Recv.List[i].Names[0])
			}
			p := variableFromExpr(p, fl, n, f.Recv.List[i].Type)
			res.Receiver = p
		}
	}

	// Change the name of the function to receiver.func
	if res.Receiver != nil {
		tmp := res.Receiver.Type
		if _, ok := tmp.(*StarType); ok {
			tmp = tmp.(*StarType).Target
		}

		res.Name = tmp.(*IdentType).Ident + "." + res.Name
	}

	res.Type = &FuncType{
		Parameters: getVariableList(p, fl, f.Type.Params),
		Results:    getVariableList(p, fl, f.Type.Results),
	}

	return res
}
