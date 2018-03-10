package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var chanCode = `
package test

type CHAN chan int

type CHAN2 chan <- int

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
		So(p.Bind(), ShouldBeNil)
		Convey("chan functions", func() {
			t, err := p.FindType("CHAN")
			So(err, ShouldBeNil)
			So(t.Equal(t), ShouldBeTrue)
			So(t.Type.Equal(t.Type), ShouldBeTrue)
			So(t.Type, ShouldHaveSameTypeAs, &ChannelType{})
			So(t.Type.Package().Equal(p), ShouldBeTrue)
			So(t.Type.String(), ShouldEqual, "chan int")
		})

		Convey("Equality", func() {
			t, err := p.FindType("CHAN")
			So(err, ShouldBeNil)
			non, err := p.FindVariable("NON")
			So(err, ShouldBeNil)
			So(non.Type.Equal(t.Type), ShouldBeTrue)
			nonc, err := p.FindVariable("NONC")
			So(err, ShouldBeNil)
			So(nonc.Type.Equal(t.Type), ShouldBeFalse)
			t2, err := p.FindType("CHAN2")
			So(err, ShouldBeNil)
			So(t.Type.Equal(t2.Type), ShouldBeFalse)
			So(t.Type.Equal(t.Type.(*ChannelType).Type), ShouldBeFalse)
		})
	})
}
