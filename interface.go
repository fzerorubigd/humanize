package humanize

import "go/ast"

// InterfaceType is the
type InterfaceType struct {
	pkg       *Package
	Functions []*Function
	Embeds    []Type // IdentType or SelectorType
}

func (i *InterfaceType) String() string {
	if len(i.Embeds) == 0 && len(i.Functions) == 0 {
		return "interface{}"
	}

	res := "interface{\n"
	for e := range i.Embeds {
		res += "\t" + i.Embeds[e].String() + "\n"
	}
	for f := range i.Functions {
		res += "\t" + i.Functions[f].Type.getDefinitionWithName(i.Functions[f].Name) + "\n"
	}
	return res + "}"
}

func (i *InterfaceType) Equal(t Type) bool {
	v, ok := t.(*InterfaceType)
	if !ok {
		return false
	}

	if !i.pkg.Equal(v.pkg) {
		return false
	}

embedLoop:
	for e := range i.Embeds {
		for ee := range v.Embeds {
			if i.Embeds[e].Equal(v.Embeds[ee]) {
				continue embedLoop
			}
		}
		return false
	}

methodLoop:
	for e := range i.Functions {
		for ee := range v.Functions {
			if i.Functions[e].Equal(v.Functions[ee]) {
				continue methodLoop
			}
		}
		return false
	}

	return true
}

func (i *InterfaceType) Package() *Package {
	return i.pkg
}

func (i *InterfaceType) lateBind() error {
	for f := range i.Functions {
		if err := i.Functions[f].lateBind(); err != nil {
			return err
		}
	}

	for f := range i.Embeds {
		if err := lateBind(i.Embeds[f]); err != nil {
			return err
		}
	}

	return nil
}

func getInterface(p *Package, f *File, t *ast.InterfaceType) Type {
	// TODO : interface may refer to itself I need more time to implement this
	iface := &InterfaceType{}
	for i := range t.Methods.List {
		res := Function{}
		// The method name is mandatory and always 1
		if len(t.Methods.List[i].Names) > 0 {
			res.Name = nameFromIdent(t.Methods.List[i].Names[0])

			res.Docs = docsFromNodeDoc(t.Methods.List[i].Doc)
			typ := newType(p, f, t.Methods.List[i].Type)
			res.Type = typ.(*FuncType)
			iface.Functions = append(iface.Functions, &res)
		} else {
			// This is the embeded interface
			embed := newType(p, f, t.Methods.List[i].Type)
			iface.Embeds = append(iface.Embeds, embed)
		}

	}
	return iface
}
