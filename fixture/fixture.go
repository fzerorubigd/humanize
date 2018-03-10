// +build ignore

package fixture

import (
	"os"

	"github.com/fzerorubigd/onion"
)

var x, err = os.Open("the_file")

var y *os.File

type f struct {
}

type f2 struct {
	*f
}

type f3 struct {
	*f
	*onion.Onion
}

// T1 test
type T1 interface {
	Test()
}

// T2 Test
type T2 interface {
	TestStar()
}

// T3 Test
type T3 interface {
	Test()
	TestStar()
}

// T4 test
type T4 interface {
	T3
	SetDelimiter(d string)
}

// NewF test
func NewF() *f {
	return &f{}
}

// NewFile test
func NewFile() *os.File {
	return y
}

// NoReturn test
func NoReturn() {

}

func (f) Test() {

}

func (*f) TestStar() {

}
