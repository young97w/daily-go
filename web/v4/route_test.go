package v2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	//测试覆盖点
	//1 比较每个http方法的tree
	//2 比较每个tree下的node
	//3 比较每个node的下的path children handlefunc

	mockHandler := func(ctx *Context) {}

	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		//{
		//	method: http.MethodGet,
		//	path:   "/param/:id",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/param/:id/detail",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/param/:id/*",
		//},
	}

	r := newRouter()

	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
						handler: mockHandler,
					},
					"order": &node{
						path: "order",
						starChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
					//"param": &node{
					//	path: "param",
					//	paramChild: &node{
					//		path: ":id",
					//		children: map[string]*node{
					//			"detail": &node{
					//				path:    "detail",
					//				handler: mockHandler,
					//			},
					//		},
					//		handler: mockHandler,
					//		starChild: &node{
					//			path:    "*",
					//			handler: mockHandler,
					//		},
					//	},
					//},
				},
				handler: mockHandler,
				//根节点的*
				starChild: &node{
					path:    "*",
					handler: mockHandler,
					children: map[string]*node{
						"abc": &node{
							path:    "abc",
							handler: mockHandler,
							starChild: &node{
								path:    "*",
								handler: mockHandler,
							},
						},
					},
					//第二个*
					starChild: &node{
						path:    "*",
						handler: mockHandler,
					},
				},
			},
		},
	}

	//比较树
	msg, ok := r.equal(wantRouter)

	assert.True(t, ok, msg)

	//非法用例
	r = newRouter()

	//空字符串
	assert.PanicsWithValue(t, "路由不能为空", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	//前导没有
	assert.PanicsWithValue(t, "路由必须以/开头", func() {
		r.addRoute(http.MethodGet, "a/bc", mockHandler)
	})

	//后缀有/
	assert.PanicsWithValue(t, "路由不能以/结尾", func() {
		r.addRoute(http.MethodGet, "/bc/", mockHandler)
	})

	//根节点重复注册
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "根节点路由冲突", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})

	//普通节点重复注册
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "路由冲突，[/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})

	//使用多个////
	assert.PanicsWithValue(t, "非法路由，不允许使用//a/b, /a//b 之类的路由, [/a//b]", func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})

	// 同时注册通配符路由和参数路由
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/*", mockHandler)
		r.addRoute(http.MethodGet, "/:id", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/:id", mockHandler)
		r.addRoute(http.MethodGet, "/*", mockHandler)
	})

	// 参数冲突
	assert.PanicsWithValue(t, "web: 路由冲突，参数路由冲突，已有 :id，新注册 :name", func() {
		r.addRoute(http.MethodGet, "/a/b/c/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/c/:name", mockHandler)
	})

}

func TestFindRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}

	mockHandler := func(ctx *Context) {}

	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchInfo
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "user",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "no handler",
			method: http.MethodPost,
			path:   "/order",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "create",
					handler: mockHandler,
				},
			},
		},

		// 通配符匹配
		{
			// 命中/order/*
			name:   "star match",
			method: http.MethodPost,
			path:   "/order/delete",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:   "star in middle",
			method: http.MethodGet,
			path:   "/user/Tom/home",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/param/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/*
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/abc",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},

		{
			// 命中 /param/:id/detail
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/detail",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
	}

	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.True(t, tc.found == found)
			if !found {
				return
			}
			wantVal := reflect.ValueOf(tc.mi.n.handler)
			nVal := reflect.ValueOf(mi.n.handler)
			assert.Equal(t, wantVal, nVal)
		})
	}
}

//equal,返回不相等的msg和布尔值
func (r *router) equal(y *router) (string, bool) {
	if len(r.trees) != len(y.trees) {
		return fmt.Sprintf("路由树个数不相等"), false
	}
	//循环want route跟r(addRoute构建的)进行比较
	for method, yt := range r.trees {
		tree, ok := y.trees[method]
		if !ok {
			return fmt.Sprintf("目标树无http方法[%s]", method), false
		}

		msg, ok := yt.equal(tree)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func (n *node) equal(yn *node) (string, bool) {
	if n.path != yn.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(yn.handler)
	if nhv != yhv {
		return fmt.Sprintf("handler不相等"), false
	}

	if len(n.children) != len(yn.children) {
		return fmt.Sprintf("字节点数量不相等"), false
	}

	//比较starChild
	if n.path == "*" {
		msg, ok := n.starChild.equal(yn.starChild)
		if !ok {
			return fmt.Sprintf("通配符节点不匹配,[%s]", msg), false
		}
	}

	for path, node := range n.children {
		dst, ok := yn.children[path]
		if !ok {
			return fmt.Sprintf("目标路由缺少路径:[%s]", path), false
		}
		msg, ok := node.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}
