package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	//先校验path
	if len(path) == 0 {
		panic("web: 路由是空字符串")
	}

	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	if path == "/" {
		if r.trees[method] != nil {
			panic(fmt.Sprintf("web: 路由冲突[%s]", path))
		} else {
			r.trees[method] = &node{
				path:    "/",
				handler: handler,
			}
			return
		}
	}

	//find treepath = {string} "/order/detail"
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	//分割path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(seg)
	}

	if root.handler != nil {
		switch root.typ {
		case nodeTypeParam:
			panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", root.path, segs[len(segs)-1]))
		default:
			panic(fmt.Sprintf("web: 路由冲突[%s]", path))
		}
	}
	root.handler = handler
}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	mi := &matchInfo{n: root}
	if path == "/" {
		return mi, mi.n != nil
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, seg := range segs {
		var param string
		var child *node
		child, param, ok = root.childOf(seg)
		if !ok {
			if root.typ == nodeTypeAny {
				return mi, root.handler != nil
			} else {
				return nil, false
			}
		}
		if param != "" {
			mi.addValue(root.paramName, param)
		}
		root = child
		mi.n = root
	}

	mi.n = root
	return mi, mi.n.handler != nil
}

type nodeType int

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 正则匹配，形式 :param_name(reg_expr)
// 3. 路径参数匹配：形式 :param_name
// 4. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	typ nodeType

	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp
}

// child 返回子节点
// 第一个返回值 *node 是命中的节点
//第二个返回值 string 代表获取到的参数
// 第三个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, string, bool) {
	root, ok := n.children[path]
	if ok {
		return root, "", root != nil
	}

	//匹配正则，路径参数，通配符
	switch {
	case n.regChild != nil:
		regStr := n.regChild.regExpr.FindString(path)
		return n.regChild, regStr, regStr != ""
	case n.paramChild != nil:
		return n.paramChild, path, true
	case n.starChild != nil:
		return n.starChild, "", true
	}

	//所有都未命中
	return nil, "", false
}

// childOrCreate 查找子节点，
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(path string) *node {
	//通配符
	if path == "*" {
		//检测参数路由
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		//检测正则
		if n.regChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", path))
		}

		if n.starChild == nil {
			n.starChild = &node{
				typ:  nodeTypeAny,
				path: path,
			}
		}
		return n.starChild
	}

	//路径参数路由
	if path[0] == ':' {
		//arr 为 [path ，paramName，正则表达式]
		arr := regexp.MustCompile(`^:(\w+)\((.*)\)$`).FindStringSubmatch(path)

		//路径参数路由
		if len(arr) == 0 {
			if n.starChild != nil {
				panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
			}

			if n.regChild != nil {
				panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
			}

			if n.paramChild == nil {
				n.paramChild = &node{
					typ:       nodeTypeParam,
					path:      path,
					paramName: path[1:],
				}
			}
			return n.paramChild
		}

		//正则路由
		//空正则表达式
		if arr[2] == "" {
			panic(fmt.Sprintf("web: 空正则表达式。[%s]", path))
		}

		//通配符
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
		}
		//路径参数
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [%s]", path))
		}
		if n.regChild == nil {
			n.regChild = &node{
				typ:       nodeTypeReg,
				path:      path,
				paramName: arr[1],
				regExpr:   regexp.MustCompile(arr[2]),
			}
		}
		return n.regChild
	}

	//静态路由
	if n.children == nil {
		n.children = map[string]*node{}
	}

	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
