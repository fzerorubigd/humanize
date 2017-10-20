package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var vr = `
package test

var i int

var i1 = 10

var i2 = 10.6

var i3 = "test"

var i4 = 'C'

var i5 = []int{1,2,3,4,5}

var i6 = 10i

`

func TestVariable(t *testing.T) {
	Convey("Variable test", t, func() {
		var p = &Package{}

		f, err := ParseFile(vr, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		Convey("Normal define", func() {
			i, err := p.FindVariable("i")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "int")
		})

		Convey("by value i1", func() {
			i, err := p.FindVariable("i1")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i1")
			So(i.Name, ShouldEqual, "i1")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "int")
		})

		Convey("by value i2", func() {
			i, err := p.FindVariable("i2")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i2")
			So(i.Name, ShouldEqual, "i2")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "float64")
		})

		Convey("by value i3", func() {
			i, err := p.FindVariable("i3")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i3")
			So(i.Name, ShouldEqual, "i3")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "string")
		})

		Convey("by value i4", func() {
			i, err := p.FindVariable("i4")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i4")
			So(i.Name, ShouldEqual, "i4")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "char")
		})

		Convey("by value i5", func() {
			i, err := p.FindVariable("i5")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i5")
			So(i.Name, ShouldEqual, "i5")
			So(i.Type.(*ArrayType).Type.(*IdentType).Ident, ShouldEqual, "int")
		})

		Convey("by value i6", func() {
			i, err := p.FindVariable("i6")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i6")
			So(i.Name, ShouldEqual, "i6")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "complex64")
		})
	})
}
