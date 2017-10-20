package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var fl = `
package test

WRONG!

`

func TestWrongFile(t *testing.T) {
	Convey("Wrong file", t, func() {
		_, err := ParseFile(fl, &Package{})
		So(err, ShouldNotBeNil)
	})
}
