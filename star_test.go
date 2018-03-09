package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var st = `
package test

type STAR *int

var x *int

var y *int

`

func TestStartType(t *testing.T) {
	Convey("Star test", t, func() {
		var p = &Package{}

		f, err := ParseFile(st, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		Convey("Normal define", func() {
			i, err := p.FindType("STAR")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "STAR")
			So(i.Type, ShouldHaveSameTypeAs, &StarType{})
			So(i.Type.(*StarType).Target.(*IdentType).Ident, ShouldEqual, "int")
			So(i.Type.Package(), ShouldEqual, p)
		})

		Convey("Equality", func() {
			So(p.Bind(), ShouldBeNil)
			x, err := p.FindVariable("x")
			So(err, ShouldBeNil)
			So(x.Type, ShouldHaveSameTypeAs, &StarType{})
			y, err := p.FindVariable("y")
			So(err, ShouldBeNil)

			So(x.Type.Equal(y.Type), ShouldBeTrue)
		})

	})
}
