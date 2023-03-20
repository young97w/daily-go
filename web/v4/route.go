package v2

import (
	"fmt"
	"strings"
)

type router struct {
	//按照http方法组织的
	//GET POST PUT DELETE等
	trees map[string]*node
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc

	//通配符匹配
	starChild *node

	//路径参数匹配
	paramChild *node
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (mi *matchInfo) addValue(key, value string) {
	if mi.pathParams == nil {
		mi.pathParams = map[string]string{key: value}
	}
	mi.pathParams[key] = value
}

func newRouter() router {
	return router{map[string]*node{}}
}

func (r *router) addRoute(method, path string, handler HandleFunc) {
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
	root, ok := r.trees[method]

	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	//返回
	if path == "/" {
		if root.handler != nil {
			panic("根节点路由冲突")
		}
		root.handler = handler
		return
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

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	//find tree
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	//校验path
	if len(path) == 0 {
		return nil, false
	}

	if path[0] != '/' {
		return nil, false
	}

	//末尾不能为/结尾
	if path != "/" && path[len(path)-1] == '/' {
		return nil, false
	}

	mi := &matchInfo{n: root}

	//如果是根节点
	if path == "/" {
		return mi, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, seg := range segs {
		var paramMatch bool
		mi.n, paramMatch, ok = mi.n.ChildOf(seg)
		if !ok {
			return nil, false
		}
		if paramMatch {
			mi.addValue(mi.n.path[1:], seg)
		}
	}
	return mi, true
}

// child 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否是命中参数路由
// 第三个返回值 bool 代表是否命中
func (n *node) ChildOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		//可能有param节点
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		//可能有通配符匹配
		return n.starChild, false, n.starChild != nil

	}
	return child, false, ok
}

func (n *node) ChildOrCreate(path string) *node {
	if path == "*" {
		//参数路径冲突
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}

		if n.starChild == nil {
			n.starChild = &node{
				path: "*",
			}
		}
		return n.starChild
	}

	//路径参数匹配
	if path[0] == ':' {
		//通用匹配冲突
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}

		//
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
			}
		} else {
			n.paramChild = &node{
				path: path,
			}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		//创建child
		child = &node{
			path: path,
		}
		//将child挂在父节点上
		n.children[path] = child
	}
	return child
}
