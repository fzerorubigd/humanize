package humanize

import (
	"go/ast"
	"strings"
)

type Docs []string

func (d Docs) String() string {
	return strings.Join(d, "\n")
}

func docsFromNodeDoc(cgs ...*ast.CommentGroup) Docs {
	var res Docs
	for _, cg := range cgs {
		if cg != nil {
			for i := range cg.List {
				res = append(res, cg.List[i].Text)
			}
		}
	}
	return res
}
