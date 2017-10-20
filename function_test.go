package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var fn = `
package test

// NoParam
// NoParam line 2
func NoParam() {

}

// NoReturn line 1
func NoReturn(a int, b string, c,d int64, _ bool) {

}

// OneReturn first
func OneReturn() error {
     return nil
}

// MultiReturn function
func MultiReturn(a,b,c int64) (d int64,e string,f error) {
      return 0, "", nil
}

type Boogh string

func (b Boogh) String() string {
       return string(b)
}

func (x *Boogh) String2() string {
return ""
}

func (Boogh) Some() {

}

func Some(int) {

}

func (Boogh) Some1(int) {

}

func Some1(string) {

}

func (Boogh) Some2() int {
	return 0
}

func Some2() error {
	return nil
}


`

func TestFunctionData(t *testing.T) {
	Convey("Function parser test", t, func() {
		var p = &Package{}
		f, err := ParseFile(fn, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)
		So(p.Bind(), ShouldBeNil)
		Convey("Fact about NoParam", func() {
			fn, err := p.FindFunction("NoParam")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "NoParam")
			So(len(fn.Type.Parameters), ShouldEqual, 0)
			So(len(fn.Type.Results), ShouldEqual, 0)
			So(fn.Receiver, ShouldBeNil)
			So(fn.Docs[0], ShouldEqual, "// NoParam")
			So(fn.Docs[1], ShouldEqual, "// NoParam line 2")
		})

		Convey("Fact about NoReturn", func() {
			fn, err := p.FindFunction("NoReturn")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "NoReturn")
			So(fn.Receiver, ShouldBeNil)
			So(len(fn.Type.Parameters), ShouldEqual, 5)
			So(len(fn.Type.Results), ShouldEqual, 0)
			So(fn.Type.Parameters[0].Name, ShouldEqual, "a")
			So(fn.Type.Parameters[0].Type.(*IdentType).Ident, ShouldEqual, "int")
			So(fn.Type.Parameters[1].Name, ShouldEqual, "b")
			So(fn.Type.Parameters[1].Type.(*IdentType).Ident, ShouldEqual, "string")
			So(fn.Type.Parameters[2].Name, ShouldEqual, "c")
			So(fn.Type.Parameters[2].Type.(*IdentType).Ident, ShouldEqual, "int64")
			So(fn.Type.Parameters[3].Name, ShouldEqual, "d")
			So(fn.Type.Parameters[3].Type.(*IdentType).Ident, ShouldEqual, "int64")
			So(fn.Type.Parameters[4].Name, ShouldEqual, "_")
			So(fn.Type.Parameters[4].Type.(*IdentType).Ident, ShouldEqual, "bool")
		})

		Convey("Fact about OneReturn", func() {
			fn, err := p.FindFunction("OneReturn")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "OneReturn")
			So(fn.Receiver, ShouldBeNil)
			So(len(fn.Type.Parameters), ShouldEqual, 0)
			So(len(fn.Type.Results), ShouldEqual, 1)
			So(fn.Type.Results[0].Name, ShouldBeEmpty)
			So(fn.Type.Results[0].Type.(*IdentType).Ident, ShouldEqual, "error")
		})

		Convey("Fact about MultiReturn", func() {
			fn, err := p.FindFunction("MultiReturn")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "MultiReturn")
			So(fn.Receiver, ShouldBeNil)
			So(len(fn.Type.Parameters), ShouldEqual, 3)
			So(len(fn.Type.Results), ShouldEqual, 3)
			So(fn.Type.Parameters[0].Name, ShouldEqual, "a")
			So(fn.Type.Parameters[0].Type.(*IdentType).Ident, ShouldEqual, "int64")
			So(fn.Type.Parameters[1].Name, ShouldEqual, "b")
			So(fn.Type.Parameters[1].Type.(*IdentType).Ident, ShouldEqual, "int64")
			So(fn.Type.Parameters[2].Name, ShouldEqual, "c")
			So(fn.Type.Parameters[2].Type.(*IdentType).Ident, ShouldEqual, "int64")

			So(fn.Type.Results[0].Name, ShouldEqual, "d")
			So(fn.Type.Results[0].Type.(*IdentType).Ident, ShouldEqual, "int64")
			So(fn.Type.Results[1].Name, ShouldEqual, "e")
			So(fn.Type.Results[1].Type.(*IdentType).Ident, ShouldEqual, "string")
			So(fn.Type.Results[2].Name, ShouldEqual, "f")
			So(fn.Type.Results[2].Type.(*IdentType).Ident, ShouldEqual, "error")
		})

		Convey("Fact about String", func() {
			fn, err := p.FindFunction("Boogh.String")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "Boogh.String")
			So(fn.Receiver.Type.(*IdentType).Ident, ShouldEqual, "Boogh")
			So(fn.Receiver.Name, ShouldEqual, "b")
			So(len(fn.Type.Parameters), ShouldEqual, 0)
			So(len(fn.Type.Results), ShouldEqual, 1)
			So(fn.Type.Results[0].Name, ShouldBeEmpty)
			So(fn.Type.Results[0].Type.(*IdentType).Ident, ShouldEqual, "string")
		})

		Convey("Fact about String2", func() {
			fn, err := p.FindFunction("Boogh.String2")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "Boogh.String2")
			So(fn.Receiver.Type.(*StarType).Target.(*IdentType).Ident, ShouldEqual, "Boogh")
			So(fn.Receiver.Name, ShouldEqual, "x")
			So(len(fn.Type.Parameters), ShouldEqual, 0)
			So(len(fn.Type.Results), ShouldEqual, 1)
			So(fn.Type.Results[0].Name, ShouldBeEmpty)
			So(fn.Type.Results[0].Type.(*IdentType).Ident, ShouldEqual, "string")
		})

		Convey("Fact about Some", func() {
			fn, err := p.FindFunction("Boogh.Some")
			So(err, ShouldBeNil)
			So(fn.Name, ShouldEqual, "Boogh.Some")
			So(fn.Receiver.Type.(*IdentType).Ident, ShouldEqual, "Boogh")
			So(len(fn.Type.Parameters), ShouldEqual, 0)
			So(len(fn.Type.Results), ShouldEqual, 0)
		})

		Convey("compare functions", func() {
			_, err := p.FindFunction("Some")
			So(err, ShouldBeNil)
			_, err = p.FindFunction("Boogh.Some")
			So(err, ShouldBeNil)

			//So(compareFunc(fn1, fn2), ShouldBeFalse)

			_, err = p.FindFunction("Some1")
			So(err, ShouldBeNil)
			_, err = p.FindFunction("Boogh.Some1")
			So(err, ShouldBeNil)

			//So(compareFunc(fn1, fn2), ShouldBeFalse)

			_, err = p.FindFunction("Some2")
			So(err, ShouldBeNil)
			_, err = p.FindFunction("Boogh.Some2")
			So(err, ShouldBeNil)

			//So(compareFunc(fn1, fn2), ShouldBeFalse)
		})
	})
}
