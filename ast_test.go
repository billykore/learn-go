package learning

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAST(t *testing.T) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "test.go", nil, parser.ParseComments)
	assert.NoError(t, err)

	t.Log("doc", f.Doc)
	t.Log("package", f.Package)
	t.Log("name", f.Name)

	t.Log("fileStart", f.FileStart)
	t.Log("fileEnd", f.FileEnd)

	if f.Imports != nil {
		for i, im := range f.Imports {
			t.Log("import", i, im.Path.Value)
		}
	}

	if f.Comments != nil {
		for i, c := range f.Comments {
			t.Log("comment", i, c.Text())
		}
	}
}
