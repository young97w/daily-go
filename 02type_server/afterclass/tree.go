package afterclass

import "fmt"

type Tree interface {
	GetValue() int
	SetValue(int)
	GetParent() Tree
	SetParent(Tree)
	GetChildren() []Tree
	AddChild(Tree)
	RemoveChild(Tree)
}

// 二叉树
type binaryTree struct {
	value  int
	parent *binaryTree
	left   *binaryTree
	right  *binaryTree
}

// 多叉树
type mutliWayTree struct {
	value    int
	parent   *binaryTree
	children []*mutliWayTree
}

func (bt *binaryTree) GetValue() int {
	return bt.value
}

func (bt *binaryTree) SetValue(val int) {
	bt.value = val
}

func (bt *binaryTree) GetParent() Tree {
	return bt.parent
}

func (bt *binaryTree) SetParent(p Tree) {
	bt.parent = p.(*binaryTree)
}

func (bt *binaryTree) GetChildren() []Tree {
	children := []Tree{}
	if bt.left != nil {
		children = append(children, bt.left)
	}
	if bt.right != nil {
		children = append(children, bt.right)
	}
	return children
}

func (bt *binaryTree) AddChild(child Tree) {
	if bt.left == nil {
		bt.left = child.(*binaryTree)
	} else if bt.right == nil {
		bt.right = child.(*binaryTree)
	}
	child.SetParent(bt)
}

func (bt *binaryTree) RemoveChild(child Tree) {
	if bt.left == child.(*binaryTree) {
		bt.left = nil
	} else if bt.right == child.(*binaryTree) {
		bt.right = nil
	}
	child.SetParent(nil)
}

func levelOrderTraversal(root *binaryTree) {
	if root == nil {
		return
	}

	queue := []*binaryTree{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Printf("%d ", node.value)

		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}
}
