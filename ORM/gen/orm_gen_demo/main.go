package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type OrmFile struct {
	File
	Ops []string
}

func main() {

}

func gen(w io.Writer, srcFile string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	tv := &SingleFileEntryVisitor{}
	ast.Walk(tv, f)
	file := tv.Get()
	fmt.Println(file)
	return nil
}
