package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var identCode = `
package test

import "net/http"

var x http.Dir

func Test(xx http.Dir) {

}

`

func TestSelector(t *testing.T) {
	Convey("array test", t, func() {
		var p = &Package{
			Name: "test",
		}
		f, err := ParseFile(identCode, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		So(p.Bind(), ShouldBeNil)
		Convey("test ident parsing", func() {
			x, err := p.FindVariable("x")
			So(err, ShouldBeNil)
			So(x.Type, ShouldHaveSameTypeAs, &SelectorType{})

			fn, err := p.FindFunction("Test")
			So(err, ShouldBeNil)
			t := fn.Type.Parameters[0]

			So(t.Type, ShouldHaveSameTypeAs, &SelectorType{})
			So(t.Type.Package().Name, ShouldEqual, "http")
		})
	})
}
