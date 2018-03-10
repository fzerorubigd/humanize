package humanize

import (
	"go/ast"
)

// ChannelType is the channel type in go source code
type ChannelType struct {
	pkg *Package

	Direction ast.ChanDir
	Type      Type
}

func (c *ChannelType) lateBind() error {
	return lateBind(c.Type)
}

// String represent string version of the data
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

// Package is the package of channel
func (c *ChannelType) Package() *Package {
	return c.pkg
}

// Equal check if two channel are exported
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
