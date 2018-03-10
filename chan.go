package humanize

import (
	"go/ast"
)

type ChannelType struct {
	pkg *Package

	Direction ast.ChanDir
	Type      Type
}

func (c *ChannelType) lateBind() error {
	return lateBind(c.Type)
}

func (c *ChannelType) String() string {
	switch c.Direction {
	case ast.SEND:
		return "chan<- " + c.Type.String()
	case ast.RECV:
		return "<-chan " + c.Type.String()
	default:
		return "chan " + c.Type.String()
	}
}

func (c *ChannelType) Package() *Package {
	return c.pkg
}

func (c *ChannelType) Equal(t Type) bool {
	v, ok := t.(*ChannelType)
	if !ok {
		return false
	}

	if c.Direction != v.Direction {
		return false
	}

	return c.Type.Equal(v.Type)
}

func getChannel(p *Package, f *File, t *ast.ChanType) Type {
	return &ChannelType{
		pkg:       p,
		Direction: t.Dir,
		Type:      newType(p, f, t.Value),
	}
}
