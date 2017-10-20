package main

import (
	"log"

	"github.com/fzerorubigd/humanize"
	"github.com/kr/pretty"
)

func main() {
	p, err := humanize.ParsePackage("github.com/fzerorubigd/humanize/cmd/test/fix")
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Bind(); err != nil {
		log.Fatal(err)
	}
	pretty.Print(p)
}
