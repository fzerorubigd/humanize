package humanize

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var imprt1 = `
package test

// Doc
import (
	"testing"
	onn "github.com/fzerorubigd/humanize"
	_ "github.com/fzerorubigd/humanize/fixture/imprt"
	// Other
	. "github.com/smartystreets/goconvey/convey"
	"github.com/fzerorubigd/humanize/fixture/imprt"
	// Not a valid go package
	"github.com/fzerorubigd/humanize/fixture/notvalid"

	"invalid/package/name/haha"
)

`

func TestImport(t *testing.T) {
	Convey("Import test ", t, func() {
		var p = &Package{}
		f, err := ParseFile(imprt1, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)
		So(p.Bind(), ShouldBeNil)
		Convey("Import testing", func() {
			i, err := p.FindImport("testing")
			So(err, ShouldBeNil)
			So(i.Package, ShouldEqual, "testing")
			So(i.Canonical, ShouldEqual, "testing")
			So(i.Path, ShouldEqual, "testing")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
			So(i.String(), ShouldEqual, `"testing"`)
		})

		Convey("Import canonical", func() {
			i, err := p.FindImport("onn")
			So(err, ShouldBeNil)
			i2, err := p.FindImport("github.com/fzerorubigd/humanize")
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(i, i2), ShouldBeTrue)
			So(i.Package, ShouldEqual, "humanize")
			So(i.Canonical, ShouldEqual, "onn")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/humanize")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")

			pk, err := i.LoadPackage()
			So(err, ShouldBeNil)
			So(pk.Path, ShouldEqual, "github.com/fzerorubigd/humanize")
			So(i.String(), ShouldEqual, `onn "github.com/fzerorubigd/humanize"`)
		})

		Convey("Import pq", func() {
			i, err := p.FindImport("github.com/fzerorubigd/humanize/fixture/imprt")
			So(err, ShouldBeNil)
			_, err = p.FindImport("_")
			So(err, ShouldNotBeNil)
			So(i.Canonical, ShouldEqual, "_")
			So(i.Package, ShouldEqual, "imprt")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/humanize/fixture/imprt")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
		})

		Convey("Import convey", func() {
			i, err := p.FindImport("github.com/smartystreets/goconvey/convey")
			So(err, ShouldBeNil)
			So(i.Canonical, ShouldEqual, ".")
			So(i.Package, ShouldEqual, "convey")
			So(i.Path, ShouldEqual, "github.com/smartystreets/goconvey/convey")
			So(len(i.Docs), ShouldEqual, 2)
			So(i.Docs[0], ShouldEqual, "// Doc")
			So(i.Docs[1], ShouldEqual, "// Other")
		})

		Convey("Import test", func() {
			i, err := p.FindImport("imprt")
			So(err, ShouldBeNil)
			i2, err := p.FindImport("github.com/fzerorubigd/humanize/fixture/imprt")
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(i, i2), ShouldBeTrue)
			So(i.Package, ShouldEqual, "imprt")
			So(i.Canonical, ShouldEqual, "_")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/humanize/fixture/imprt")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
			So(i.Folder, ShouldEndWith, "github.com/fzerorubigd/humanize/fixture/imprt")

		})

		Convey("Import invalid folder", func() {
			i, err := p.FindImport("github.com/fzerorubigd/humanize/fixture/notvalid")
			So(err, ShouldBeNil)
			So(i.Package, ShouldEqual, "notvalid")
			So(i.Canonical, ShouldEqual, "notvalid")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/humanize/fixture/notvalid")
			So(i.Folder, ShouldEqual, "")
			pk, err := i.LoadPackage()
			So(err, ShouldNotBeNil)
			So(pk, ShouldBeNil)
		})

		Convey("Import not exist folder", func() {
			i, err := p.FindImport("invalid/package/name/haha")
			So(err, ShouldBeNil)
			So(i.Package, ShouldEqual, "haha")
			So(i.Canonical, ShouldEqual, "haha")
			So(i.Path, ShouldEqual, "invalid/package/name/haha")
			So(i.Folder, ShouldEqual, "")
		})

		Convey("Import incorrect", func() {
			i, err := p.FindImport("i/am/not/in/the/source")
			So(err, ShouldNotBeNil)
			So(i, ShouldBeNil)
		})
	})
}
