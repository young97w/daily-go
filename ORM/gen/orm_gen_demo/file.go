package main

import "go/ast"

//SingleFileEntryVisitor，代表单个 Go 文件的访问者；
//FileVisitor，代表单个 Go 文件中的结构体的访问者。
//File，代表单个 Go 文件的结构体；

// File 最终返回的结构体
type File struct {
	Package string
	Imports []string
	Types   []Type
}

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

type SingleFileEntryVisitor struct {
	file *FileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	if s.file != nil {
		return s.file.Get()
	}
	return File{}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fn, ok := node.(*ast.File)
	if !ok {
		return s
	}
	//fn.Name: package name
	s.file = &FileVisitor{Package: fn.Name.String()}
	return s.file
}

type FileVisitor struct {
	Package string
	Imports []string
	Types   []*TypeVisitor
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.TypeSpec:
		v := &TypeVisitor{name: n.Name.String()}
		f.Types = append(f.Types, v)
		// 调用typevisitor
		return v
	case *ast.ImportSpec:
		path := n.Path.Value
		if n.Name != nil && n.Name.String() != "" {
			path = n.Name.String() + " " + path
		}
		f.Imports = append(f.Imports, path)
	}
	return f
}

func (f *FileVisitor) Get() File {
	types := make([]Type, 0, len(f.Types))
	for _, t := range f.Types {
		types = append(types, t.Get())
	}
	return File{
		Package: f.Package,
		Imports: f.Imports,
		Types:   types,
	}
}

type TypeVisitor struct {
	name   string
	fields []Field
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	n, ok := node.(*ast.Field)
	if !ok {
		return t
	}
	var typ string
	switch nt := n.Type.(type) {
	case *ast.Ident:
		typ = nt.String()
	case *ast.StarExpr:
		switch xt := nt.X.(type) {
		case *ast.Ident:
			typ = "*" + xt.String()
		case *ast.SelectorExpr:
			typ = "*" + xt.X.(*ast.Ident).String() + "." + xt.Sel.String()
		}
	case *ast.ArrayType:
		typ = "[]byte"
	default:
		panic("不支持的类型")
	}
	for _, name := range n.Names {
		t.fields = append(t.fields, Field{
			Name: name.String(),
			Type: typ,
		})
	}
	return t
}

func (t *TypeVisitor) Get() Type {
	return Type{
		Name:   t.name,
		Fields: t.fields,
	}
}
