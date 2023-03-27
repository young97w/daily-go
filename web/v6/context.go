package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	RespStatusCode int
	RespData       []byte

	// 命中的路由
	MatchedRoute string

	// 万一将来有需求，可以考虑支持这个，但是需要复杂一点的机制
	// Body []byte 用户返回的响应
	// Err error 用户执行的 Error

	// 缓存的数据
	cacheQueryValues url.Values
}

//BindJSON 解码json
//读数据
func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

//FormValue 返回form值
//读数据，读form前要parseForm
func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

//QueryValue 返回query值
//读数据,应先把query缓存起来，应对多次读取
func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到key")}
	}
	return StringValue{val: vals[0]}
}

//PathValue 返回参数路径或者正则匹配的值
//读数据
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到key")}
	}
	return StringValue{val: val}
}

//SetCookie 设置Cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

//RespJSONOK 响应一个ok的消息
func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

//RespJSON 响应json消息
func (c *Context) RespJSON(code int, val any) error {
	c.Resp.WriteHeader(code)
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	//_, err = c.Resp.Write(data)
	c.RespStatusCode = code
	c.RespData = data
	return nil
}

//StringValue 方便链式调用
type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// 不能用泛型
// func (s StringValue) To[T any]() (T, error) {
//
// }
