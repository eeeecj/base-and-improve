package binarytree

type BinaryTreeImp interface {
	Insert(data interface{})
	Search(data interface{}) *Node
	InorderTraversal(node *Node, callback func(int))
	PreorderTraversal(node *Node, callback func(int))
	PostorderTraversal(node *Node, callback func(int))
}

func NewBinaryTree() BinaryTreeImp {
	return &BinaryTree{}
}

type Node struct {
	data   interface{}
	parent *Node
	left   *Node
	right  *Node
}
type BinaryTree struct {
	root *Node
}

func (t *BinaryTree) Insert(data interface{}) {
	if t.root == nil {
		t.root = &Node{data: data}
		return
	}
	cur := t.root
	for cur != nil {
		if data.(int) > t.root.data.(int) {
			if cur.right == nil {
				cur.right = &Node{data: data}
				return
			}
			cur = cur.right
		} else {
			if cur.left == nil {
				cur.left = &Node{data: data}
				return
			}
			cur = cur.left
		}
	}
}

func (t *BinaryTree) Search(data interface{}) *Node {
	if t.root == nil {
		return nil
	}
	cur := t.root
	for cur != nil {
		if cur.data.(int) == data.(int) {
			return cur
		} else if cur.data.(int) < data.(int) {
			cur = cur.right
		} else {
			cur = cur.left
		}
	}
	return nil
}

func (t *BinaryTree) InorderTraversal(node *Node, callback func(int)) {
	if node == nil {
		return
	}
	t.InorderTraversal(node.left, callback)
	callback(node.data.(int))
	t.InorderTraversal(node.right, callback)
}

func (t *BinaryTree) PreorderTraversal(node *Node, callback func(int)) {
	if node == nil {
		return
	}
	callback(node.data.(int))
	t.PreorderTraversal(node.left, callback)
	t.PreorderTraversal(node.right, callback)
}

func (t *BinaryTree) PostorderTraversal(node *Node, callback func(int)) {
	if node == nil {
		return
	}
	t.PostorderTraversal(node.left, callback)
	t.PostorderTraversal(node.right, callback)
	callback(node.data.(int))
}
