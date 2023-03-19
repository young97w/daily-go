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
	}

	mockHandler := func(ctx *Context) {}

	testCases := []struct {
		name     string
		method   string
		path     string
		found    bool
		wantNode *node
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
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			wantNode: &node{
				path:    "user",
				handler: mockHandler,
			},
		},
		{
			name:   "no handler",
			method: http.MethodPost,
			path:   "/order",
			found:  true,
			wantNode: &node{
				path: "order",
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			wantNode: &node{
				path:    "create",
				handler: mockHandler,
			},
		},
	}

	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.True(t, tc.found == found)
			if !found {
				return
			}
			wantVal := reflect.ValueOf(tc.wantNode.handler)
			nVal := reflect.ValueOf(n.handler)
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
