package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var chanCode = `
package test

type CHAN chan int

var NON =  make(chan int)

var NONC = make(CHAN)

`

func TestChan(t *testing.T) {
	Convey("chan test", t, func() {
		var p = &Package{
			Name: "test",
		}
		f, err := ParseFile(chanCode, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		Convey("chan functions", func() {
			t, err := p.FindType("CHAN")
			So(err, ShouldBeNil)
			So(t.Equal(t), ShouldBeTrue)
			So(t.Type.Equal(t.Type), ShouldBeTrue)
			So(t.Type, ShouldHaveSameTypeAs, &ChannelType{})
			So(t.Type.Package().Equal(p), ShouldBeTrue)
			So(t.Type.String(), ShouldEqual, "chan int")
		})
	})
}
