package ast

import (
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestPrintVisitor_Visit(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", goFile, parser.ParseComments)
	require.NoError(t, err)
	ast.Walk(&PrintVisitor{}, f)
}

var goFile = `package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type PrintVisitor struct {
}

func (p PrintVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return p
	}

	typ := reflect.TypeOf(node)
	val := reflect.ValueOf(node)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	fmt.Printf("val: %+v , typ %s \n", val, typ)
	return p
}
`
