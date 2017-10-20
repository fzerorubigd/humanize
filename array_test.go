package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var arrayCode = `
package test

import "net/http"

type ARRAY []int

var NON = [...]int{1,2,3,4}

type ANOTHER []int

type ANOTHERONE []int64

type AO [10]int

type A0 [20]int

type HH http.Request

func test(a ARRAY) {

}
`

func TestArray(t *testing.T) {
	Convey("array test", t, func() {
		var p = &Package{
			Name: "test",
		}
		f, err := ParseFile(arrayCode, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		Convey("array string", func() {
			t, err := f.FindType("ARRAY")
			So(err, ShouldBeNil)
			So(t.Equal(t), ShouldBeTrue)
			So(t.Type.Equal(t.Type), ShouldBeTrue)
			So(t.String(), ShouldEqual, "ARRAY []int")

			v, err := f.FindVariable("NON")
			So(err, ShouldBeNil)
			So(v.Type, ShouldHaveSameTypeAs, &EllipsisType{})
			So(v.Type.Equal(v.Type), ShouldBeTrue)
			So(v.Type.Equal(t.Type), ShouldBeFalse)
			So(t.Type.Equal(v.Type), ShouldBeFalse)

			So(t.Type.Package().Equal(v.Type.Package()), ShouldBeTrue)

			another, err := f.FindType("ANOTHER")
			So(err, ShouldBeNil)
			So(another.Equal(another), ShouldBeTrue)
			So(another.Type.Equal(another.Type), ShouldBeTrue)
			So(another.String(), ShouldEqual, "ANOTHER []int")
			So(t.Equal(another), ShouldBeFalse)
			So(another.Equal(t), ShouldBeFalse)
			// The types are same
			So(t.Type.Equal(another.Type), ShouldBeTrue)
			So(another.Type.Equal(t.Type), ShouldBeTrue)

			httpP, err := f.FindImport("http")
			So(err, ShouldBeNil)
			p2, err := httpP.LoadPackage()
			So(err, ShouldBeNil)
			req, err := p2.FindType("Request")
			So(t.Equal(req), ShouldBeFalse)
			So(t.Type.Equal(req.Type), ShouldBeFalse)

			another, err = f.FindType("AO")
			So(err, ShouldBeNil)
			So(another.Type, ShouldHaveSameTypeAs, t.Type)
			So(another.Equal(another), ShouldBeTrue)
			So(another.Type.Equal(another.Type), ShouldBeTrue)
			So(another.String(), ShouldEqual, "AO [10]int")
			So(t.Equal(another), ShouldBeFalse)
			So(another.Equal(t), ShouldBeFalse)
			So(t.Type.Equal(another.Type), ShouldBeFalse)
			So(another.Type.Equal(t.Type), ShouldBeFalse)

			next, err := f.FindType("A0")
			So(err, ShouldBeNil)
			So(next.Type, ShouldHaveSameTypeAs, another.Type)
			So(next.Type.Equal(another.Type), ShouldBeFalse)

			fn, err := f.FindFunction("test")
			So(err, ShouldBeNil)
			So(len(fn.Type.Parameters), ShouldEqual, 1)
			So(len(fn.Type.Results), ShouldEqual, 0)

			t2 := fn.Type.Parameters[0]
			So(t2.Type.Equal(t.Type), ShouldBeFalse)
			So(t2.String(), ShouldEqual, "a ARRAY")

			t3, err := FindTypeName(t2.Type)
			So(err, ShouldBeNil)
			So(t3.Equal(t), ShouldBeTrue)
		})
	})
}
