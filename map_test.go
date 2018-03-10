package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var mp = `
package maptest

type (
   MAP map[int]string
)

var m map[string]int

var n map[string]int
`

func TestMap(t *testing.T) {
	Convey("Map parser test", t, func() {
		var p = &Package{}
		f, err := ParseFile(mp, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)
		So(p.Bind(), ShouldBeNil)
		Convey("facts about MAP", func() {
			mp, err := p.FindType("MAP")
			So(err, ShouldBeNil)
			So(mp.Name, ShouldEqual, "MAP")
			So(mp.Type, ShouldHaveSameTypeAs, &MapType{})
			So(mp.String(), ShouldEqual, "MAP map[int]string")
			mpt := mp.Type.(*MapType)
			So(mpt.Key, ShouldHaveSameTypeAs, &IdentType{})
			So(mpt.Value, ShouldHaveSameTypeAs, &IdentType{})
			So(mpt.Package().Equal(p), ShouldBeTrue)

			m, err := p.FindVariable("m")
			So(err, ShouldBeNil)
			n, err := p.FindVariable("n")
			So(err, ShouldBeNil)
			So(m.Type.Equal(n.Type), ShouldBeTrue)
		})
	})
}
