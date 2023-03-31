package web

import (
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestFileUploader_Handle(t *testing.T) {

	//初始化file 处理
	f := &FileUploader{
		FileField:   "file",
		DstPathFunc: dstPath,
	}

	//初始化页面渲染
	tpl, err := template.ParseGlob("template_data/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}

	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{T: tpl}))
	s.Get("/upload", func(ctx *Context) {
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			t.Fatal(err)
		}
	}) //

	s.Post("/upload", f.Handle())

	err = s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

func dstPath(fh *multipart.FileHeader) string {
	path := filepath.Join("web_file/docs")
	err := os.MkdirAll(path, 0o666)
	if err != nil {
		log.Fatalln(err)
	}

	return filepath.Join(path, fh.Filename)
}

func TestFileDownloader_Handle(t *testing.T) {
	fd := &FileDownloader{Dir: filepath.Join("web_file", "docs")}

	s := NewHTTPServer()
	s.Get("/download", fd.Handle())

	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	optCache := StaticResourceWithCache(1024*1024*10, 100)
	optExt := StaticResourceWithExtension(map[string]string{
		"html": "text/html",
	})
	h := NewStaticResourceHandler(filepath.Join("web_file", "static"), optCache, optExt)

	s := NewHTTPServer()
	s.Get("/:file", h.Handle)

	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

func TestReadFile(t *testing.T) {
	//fmt.Println(os.Getwd())

	//open 只读
	f, err := os.Open("web_file/static/welcom.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 6000)
	//实际读的长度 <= n
	n, err := f.Read(data)

	fmt.Println(n, string(data))
}

func TestWriteFile(t *testing.T) {
	//os.O_TRUNC 为覆盖操作
	f, err := os.OpenFile("file_data/test.txt", os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("你好哇")

	n, err := f.Write(data)
	fmt.Println(n)
	n, err = f.WriteString("可以可以")
	fmt.Println(n)

}

func TestCreateFile(t *testing.T) {
	f, err := os.Create("file_data/create.txt")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	n, err := f.WriteString("hello")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(n)
}

func TestFilePath(t *testing.T) {
	path := "a/b.txt"
	//err := os.MkdirAll("a/b.txt", 0755)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("666")

	//if err != nil {
	//	if os.IsPermission(err) {
	//		// 更改文件夹权限
	//		err = os.Chmod("a", 0755)
	//		if err != nil {
	//			panic(err)
	//		}
	//		// 再次尝试创建文件夹
	//		err = os.MkdirAll("a/b", 0755)
	//		if err != nil {
	//			panic(err)
	//		}
	//	} else {
	//		panic(err)
	//	}
	//}
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilePathPrint(t *testing.T) {
	s := filepath.Clean("/../../file.txt")
	fmt.Println(s)
}
