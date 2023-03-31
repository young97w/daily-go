package web

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"log"
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl, err := template.ParseGlob("template_data/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}

	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{T: tpl}))
	s.Get("/login", func(ctx *Context) {
		err := ctx.Render("index.gohtml", nil)
		if err != nil {
			log.Fatalln(err)
		}
	})

	err = s.Start(":8081")
	if err != nil {
		panic(err)
	}

}

type User struct {
	Name string
	T    TemplateEngine
}

func TestHelloStruct(t *testing.T) {
	tpl := template.New("hello")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}

	user := User{Name: "young"}

	bs := &bytes.Buffer{}

	err = tpl.Execute(bs, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `Hello, young`, bs.String())

}

func TestHelloMap(t *testing.T) {
	//使用map
	m := map[string]string{"Name": "Young", "age": "18"}

	bs := &bytes.Buffer{}
	tpl := template.New("hello")
	tpl, err := tpl.Parse(`Hello, {{.Name}}.You are {{.age}} now.`)
	if err != nil {
		t.Fatal(err)
	}

	err = tpl.Execute(bs, m)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `Hello, Young.You are 18 now.`, bs.String())
}

func TestHelloSlice(t *testing.T) {
	//使用slice
	m := []string{"Young", "18"}

	bs := &bytes.Buffer{}
	tpl := template.New("hello")
	tpl, err := tpl.Parse(`Hello, {{index . 0}}.You are {{index . 1}} now.`)
	if err != nil {
		t.Fatal(err)
	}

	err = tpl.Execute(bs, m)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `Hello, Young.You are 18 now.`, bs.String())
}

type FuncCall struct {
	Slice []string
}

func TestHelloFunc(t *testing.T) {
	m := []string{"Young", "118"}

	bs := &bytes.Buffer{}

	tpl := template.New("hello")
	tpl, err := tpl.Parse(`{{.Print "Young" "18"}}`)
	if err != nil {
		t.Fatal(err)
	}

	err = tpl.Execute(bs, &FuncCall{Slice: m})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `Hello, Young.You are 18 now.`, bs.String())
}

func (f FuncCall) Print(name, age string) string {
	return fmt.Sprintf("Hello, %s.You are %s now.", name, age)
}

func TestForLoop(t *testing.T) {
	M := []string{"Young", "Tom", "Eizo"}

	tpl := template.New("hello")
	// - 去除空行
	tpl, err := tpl.Parse(`
{{- range $idx,$elem := . -}}
{{- $idx -}}
{{end -}}
`)
	if err != nil {
		t.Fatal(err)
	}

	bs := &bytes.Buffer{}

	err = tpl.Execute(bs, M)

	assert.Equal(t, "012", bs.String())
}

func TestForI(t *testing.T) {
	M := make([]bool, 5)

	tpl := template.New("hello")
	// - 去除空行
	tpl, err := tpl.Parse(`
{{- range $idx,$elem := . -}}
{{- $idx -}}
{{end -}}
`)
	if err != nil {
		t.Fatal(err)
	}

	bs := &bytes.Buffer{}

	err = tpl.Execute(bs, M)

	assert.Equal(t, "01234", bs.String())
}

func TestIf(t *testing.T) {
	Age := 18

	tpl := template.New("hello")
	// - 去除空行
	tpl, err := tpl.Parse(`
{{- if and (ge . 0) (lt . 10)}}
儿童 0 <= age < 10
{{else if and (ge . 10) (lt . 18)}}
青少年 10 <= age < 18
{{else}}
成年人
{{end -}}
`)
	if err != nil {
		t.Fatal(err)
	}

	bs := &bytes.Buffer{}

	err = tpl.Execute(bs, Age)

	assert.Equal(t, "成年人", bs.String())
}
