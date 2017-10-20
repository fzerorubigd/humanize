package humanize

import (
	"fmt"
)

// FindImport try to find an import in a file
func (f *File) FindImport(pkg string) (*Import, error) {
	if pkg == "" || pkg == "_" || pkg == "." {
		return nil, fmt.Errorf("import with path _/. or empty is invalid")
	}
	for i := range f.Imports {
		if f.Imports[i].Package == pkg || f.Imports[i].Canonical == pkg || f.Imports[i].Path == pkg {
			return f.Imports[i], nil
		}
	}

	return nil, fmt.Errorf("pkg %s is not found in %s", pkg, f.FileName)
}

// FindImport try to find an import in a package
func (p *Package) FindImport(pkg string) (*Import, error) {
	for i := range p.Files {
		if i, err := p.Files[i].FindImport(pkg); err == nil {
			return i, nil
		}
	}

	return nil, fmt.Errorf("pkg %s is not found in %s", pkg, p.Name)
}

func (f *File) FindType(t string) (*TypeName, error) {
	for i := range f.Types {
		if f.Types[i].Name == t {
			return f.Types[i], nil
		}
	}
	return nil, fmt.Errorf("type %s is not found in %s", t, f.FileName)
}

func (p *Package) FindType(t string) (*TypeName, error) {
	for i := range p.Files {
		if ty, err := p.Files[i].FindType(t); err == nil {
			return ty, nil
		}
	}
	return nil, fmt.Errorf("type %s is not found in %s", t, p.Name)
}

func (f *File) FindConstant(t string) (*Constant, error) {
	for i := range f.Constants {
		if f.Constants[i].Name == t {
			return f.Constants[i], nil
		}
	}
	return nil, fmt.Errorf("const %s is not found in %s", t, f.FileName)
}

func (p *Package) FindConstant(t string) (*Constant, error) {
	for i := range p.Files {
		if ct, err := p.Files[i].FindConstant(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("const %s is not found in %s", t, p.Name)
}

func (f *File) FindFunction(t string) (*Function, error) {
	for i := range f.Functions {
		if f.Functions[i].Name == t {
			return f.Functions[i], nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, f.FileName)
}

func (p *Package) FindFunction(t string) (*Function, error) {
	for i := range p.Files {
		if ct, err := p.Files[i].FindFunction(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, p.Name)
}

func (f *File) FindVariable(t string) (*Variable, error) {
	for i := range f.Variables {
		if f.Variables[i].Name == t {
			return f.Variables[i], nil
		}
	}
	return nil, fmt.Errorf("variable %s is not found in %s", t, f.FileName)
}

func (p *Package) FindVariable(t string) (*Variable, error) {
	for i := range p.Files {
		if ct, err := p.Files[i].FindVariable(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("variable %s is not found in %s", t, p.Name)
}
