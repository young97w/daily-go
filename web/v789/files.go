package web

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//把文件处理当成 一个HandleFunc

type FileUploader struct {
	// FileField 对应文件在表单的字段名称
	FileField   string
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	//可以做额外的检测
	return func(ctx *Context) {
		formFile, header, err := ctx.Req.FormFile(f.FileField)
		defer formFile.Close()
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("文件上传失败！")
			return
		}

		dst, err := os.OpenFile(f.DstPathFunc(header), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		defer dst.Close()
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("文件上传失败！")
			return
		}
		defer dst.Close()

		_, err = io.CopyBuffer(dst, formFile, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("文件上传失败！")
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("文件上传成功！")
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		name, _ := ctx.QueryValue("file").String()
		//路径校验，不能包含相对路径
		path := filepath.Join(f.Dir, filepath.Clean(name))

		fn := filepath.Base(path)

		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

type StaticResourceHandler struct {
	dir string
	//pathPrefix              string
	extensionContentTypeMap map[string]string

	//缓存静态资源
	cache       *lru.Cache
	maxFileSize int
}

//缓存文件信息
type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

type StaticResourceHandlerOption func(h *StaticResourceHandler)

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	s := &StaticResourceHandler{
		dir: dir,
		//pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
			//"html": "text/html",
		},
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func StaticResourceWithCache(maxSize, maxCount int) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		c, err := lru.New(maxCount)
		if err != nil {
			panic(err)
		}
		h.maxFileSize = maxSize
		h.cache = c
	}
}

func StaticResourceWithExtension(exts map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, c := range exts {
			h.extensionContentTypeMap[ext] = c
		}
	}
}

func (h *StaticResourceHandler) Handle(ctx *Context) {
	name, _ := ctx.PathValue("file").String()
	path := filepath.Join(h.dir, name)

	//读缓存
	if item, ok := h.readFromCache(name); ok {
		h.writeItemAsResponse(item, ctx)
	}

	f, err := os.Open(path)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}

	ext := getFileExt(name)
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		ctx.RespStatusCode = http.StatusBadRequest
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}

	item := &fileCacheItem{
		fileName:    name,
		fileSize:    len(data),
		contentType: t,
		data:        data,
	}
	h.cacheFile(item)

	h.writeItemAsResponse(item, ctx)
}

func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize <= h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

func (h *StaticResourceHandler) readFromCache(name string) (*fileCacheItem, bool) {
	if h.cache != nil {

	}
	return nil, false
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, ctx *Context) {
	ctx.Resp.Header().Set("Content-Type", item.contentType)
	ctx.Resp.Header().Set("Content-Length", strconv.Itoa(len(item.data)))
	ctx.RespStatusCode = http.StatusOK
	ctx.RespData = item.data
}

func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}

	return strings.ToLower(name[index+1:])
}
