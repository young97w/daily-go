package v2

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	//按照http方法组织的
	//GET POST PUT DELETE等
	tress map[string]*node
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc
}

func newRoute() router {
	return router{map[string]*node{}}
}

func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (r *router) AddRoute(method, path string, handler HandleFunc) {
	if len(path) == 0 {
		panic("路由不能为空")
	}

	if path[0] != '/' {
		panic("路由必须以/开头")
	}

	//末尾不能为/结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("路由不能以/结尾")
	}

	//如果获取失败，root是nil，后续使用应该先初始化
	root, ok := r.tress[method]

	if !ok {
		root = &node{
			path: "/",
		}
		r.tress[method] = root
	}

	segs := strings.Split(path[1:], "/")

	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("非法路由，不允许使用//a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.ChildOrCreate(seg)
	}

	//遍历完之后，检查handleFunc
	if root.handler != nil {
		panic(fmt.Sprintf("路由冲突，[%s]", path))
	}

	root.handler = handler

}

func (n *node) ChildOrCreate(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[seg]
	if !ok {
		//创建child
		child = &node{
			path: seg,
		}
		//将child挂在父节点上
		n.children[seg] = child
	}
	return child
}
