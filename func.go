package humanize

import (
	"go/ast"
	"strings"
)

// FuncType is the single function
type FuncType struct {
	pkg *Package

	Parameters []*Variable
	Results    []*Variable
}

func (f *FuncType) getDefinitionWithName(name string) string {
	return "func " + name + f.Sign()
}

func (f *FuncType) Sign() string {
	var args, res []string
	for a := range f.Parameters {
		args = append(args, f.Parameters[a].Type.String())
	}

	for a := range f.Results {
		res = append(res, f.Results[a].Type.String())
	}

	result := "(" + strings.Join(args, ",") + ")"
	if len(res) > 1 {
		result += " (" + strings.Join(res, ",") + ")"
	} else {
		result += " " + strings.Join(res, ",")
	}

	return result
}

func (f *FuncType) String() string {
	return "func " + f.Sign()
}

func (f *FuncType) Package() *Package {
	return f.pkg
}

func (f *FuncType) Equal(t Type) bool {
	v, ok := t.(*FuncType)
	if !ok {
		return false
	}

	if !f.pkg.Equal(v.pkg) {
		return false
	}

	if len(f.Parameters) != len(v.Parameters) {
		return false
	}

	if len(f.Results) != len(v.Results) {
		return false
	}

	for i := range f.Parameters {
		if !f.Parameters[i].Type.Equal(v.Parameters[i].Type) {
			return false
		}
	}

	for i := range f.Results {
		if !f.Results[i].Type.Equal(v.Results[i].Type) {
			return false
		}
	}

	return true

}

func (f *FuncType) lateBind() error {
	for i := range f.Parameters {
		if err := f.Parameters[i].lateBind(); err != nil {
			return err
		}
	}
	for i := range f.Results {
		if err := f.Results[i].lateBind(); err != nil {
			return err
		}
	}
	return nil
}

func getVariableList(p *Package, fl *File, f *ast.FieldList) []*Variable {
	var res []*Variable
	if f == nil {
		return res
	}
	for i := range f.List {
		n := f.List[i].Names
		if n != nil {
			for in := range n {
				p := variableFromExpr(p, fl, nameFromIdent(n[in]), f.List[i].Type)
				res = append(res, p)
			}
		} else {
			// Its probably without name part (ie return variable)
			p := variableFromExpr(p, fl, "", f.List[i].Type)
			res = append(res, p)
		}
	}

	return res
}

func getFunc(p *Package, f *File, t *ast.FuncType) Type {
	return &FuncType{
		pkg:        p,
		Parameters: getVariableList(p, f, t.Params),
		Results:    getVariableList(p, f, t.Results),
	}
}
