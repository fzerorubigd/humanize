package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var cr = `
package test

const (
   i = iota
   j
)

const s float64 = 10

const i1 = 10

const i2 = 10.6

const i3 = "test"

const i4 = 'C'

const i6 = 10i

`

var otherC = `
package test2

const i7 = 10i

const i6 = 1000
`

func TestConstant(t *testing.T) {
	Convey("constant test", t, func() {
		f, err := ParseFile(cr, &Package{})
		So(err, ShouldBeNil)
		var p = &Package{
			Name: "test",
		}
		p.Files = append(p.Files, f)
		So(p.Bind(), ShouldBeNil)
		Convey("Normal define", func() {
			i, err := p.FindConstant("i")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "iota")

			i, err = p.FindConstant("j")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "j")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "iota")

		})
		Convey("by value s", func() {
			i, err := p.FindConstant("s")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "s")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "float64")

			So(i.String(), ShouldEqual, "s float64 = 10")
		})
		Convey("by value i1", func() {
			i, err := p.FindConstant("i1")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i1")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "int")
		})

		Convey("by value i2", func() {
			i, err := p.FindConstant("i2")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i2")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "float64")
		})

		Convey("by value i3", func() {
			i, err := p.FindConstant("i3")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i3")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "string")
		})

		Convey("by value i4", func() {
			i, err := p.FindConstant("i4")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i4")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "char")
		})

		Convey("by value i6", func() {
			i, err := p.FindConstant("i6")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "i6")
			So(i.Type.(*IdentType).Ident, ShouldEqual, "complex64")

			Convey("equality", func() {
				So(i.Equal(i), ShouldBeTrue)
				i2, err := p.FindConstant("i4")
				So(err, ShouldBeNil)
				So(i.Equal(i2), ShouldBeFalse)
			})

			Convey("other package", func() {
				f, err := ParseFile(otherC, &Package{})
				So(err, ShouldBeNil)
				var p2 = &Package{
					Name: "test2",
				}
				p2.Files = append(p2.Files, f)
				So(p2.Bind(), ShouldBeNil)

				i3, err := p2.FindConstant("i6")
				So(err, ShouldBeNil)
				So(i3.Equal(i), ShouldBeFalse)

				i4, err := p2.FindConstant("i7")
				So(err, ShouldBeNil)
				So(i4.Equal(i), ShouldBeTrue)
			})
		})
	})
}
