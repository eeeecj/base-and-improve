package AVL

type AVLNode struct {
	data        int
	left, right *AVLNode
	height      int
}
type AVLTree struct {
	root *AVLNode
}

func NewAVLTree(data int) *AVLTree {
	return &AVLTree{root: &AVLNode{data: data, height: 1}}
}

func (t *AVLTree) Find(data int) *AVLNode {
	if t.root == nil {
		return nil
	}
	return t.root.Find(data)
}

func (node *AVLNode) Find(data int) *AVLNode {
	if node == nil {
		return nil
	}
	if data == node.data {
		return node
	} else if data < node.data {
		return node.left.Find(data)
	} else {
		return node.right.Find(data)
	}
}

func (t *AVLTree) Insert(data int) {
	t.root = t.root.Insert(data)
}
func (node *AVLNode) Insert(data int) *AVLNode {
	if node == nil {
		return &AVLNode{data: data, height: 1}
	}
	if node.data == data {
		return node
	}
	var newNode *AVLNode
	if node.data > data {
		node.left = node.left.Insert(data)
		bf := node.BalanceFactor()
		if bf == 2 {
			if data < node.left.data {
				newNode = node.RightRotate()
			} else {
				newNode = node.LeftRightNode()
			}
		}
	} else {
		node.right = node.right.Insert(data)
		bf := node.BalanceFactor()
		if bf == -2 {
			if data > node.right.data {
				newNode = node.LeftRotate()
			} else {
				newNode = node.RightLeftNode()
			}
		}
	}
	if newNode == nil {
		node.UpdateHeight()
		return node
	} else {
		//newNode.UpdateHeight()
		return newNode
	}
}

func (node *AVLNode) BalanceFactor() int {
	leftHeight, rightHeight := 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	return leftHeight - rightHeight
}

func (node *AVLNode) RightRotate() *AVLNode {
	pivot := node.left
	pivotR := pivot.right
	pivot.right = node
	node.left = pivotR

	node.UpdateHeight()
	pivot.UpdateHeight()
	return pivot
}

func (node *AVLNode) LeftRotate() *AVLNode {
	pivot := node.right
	pivotL := pivot.left
	pivot.left = node
	node.right = pivotL

	node.UpdateHeight()
	pivot.UpdateHeight()
	return pivot
}

func (node *AVLNode) LeftRightNode() *AVLNode {
	node.left = node.left.LeftRotate()
	return node.RightRotate()
}

func (node *AVLNode) RightLeftNode() *AVLNode {
	node.right = node.right.RightRotate()
	return node.LeftRotate()
}

func (node *AVLNode) UpdateHeight() {
	if node == nil {
		return
	}
	leftHeight, rightHeight := 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	node.height = max(leftHeight, rightHeight) + 1
}

func (t *AVLTree) Traverse() []int {
	return t.root.Traverse()
}

//func (t *AVLTree) Traverse() {
//	t.root.Traverse()
//}

func (node *AVLNode) Traverse() []int {
	if node == nil {
		return []int{}
	}
	res := []int{}
	res = append(res, node.left.Traverse()...)
	res = append(res, node.data)
	res = append(res, node.right.Traverse()...)
	return res
}

//func (node *AVLNode) Traverse() {
//	if node == nil {
//		return
//	}
//	node.left.Traverse()
//	fmt.Printf("%d(%d:%d) ", node.data, node.BalanceFactor(), node.height)
//	node.right.Traverse()
//}

func (t *AVLTree) IsAVLTree() bool {
	if t == nil || t.root == nil {
		return true
	}
	if t.root.IsBalanced() {
		return true
	}
	return false
}
func (node *AVLNode) IsBalanced() bool {
	if node == nil {
		return true
	}
	if node.left == nil && node.right == nil {
		if node.height == 1 {
			return true
		}
		return false
	} else if node.left != nil && node.right != nil {
		if node.left.data > node.data || node.right.data < node.data {
			return false
		}
		bf := node.left.height - node.right.height
		if abs(bf) > 1 {
			return false
		}
		if node.left.height > node.right.height {
			if node.height != node.left.height+1 {
				return false
			}
		} else {
			if node.height != node.right.height+1 {
				return false
			}
		}
		if !node.left.IsBalanced() {
			return false
		}
		if !node.right.IsBalanced() {
			return false
		}
	} else {
		if node.right != nil {
			if node.right.height == 1 && node.right.left == nil && node.right.right == nil {
				if node.right.data < node.data {
					return false
				}
			} else {
				return false
			}
		} else {
			if node.left.height == 1 && node.left.left == nil && node.left.right == nil {
				if node.left.data > node.data {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

func (t *AVLTree) Delete(data int) {
	if t.root == nil {
		return
	}

}

func (node *AVLNode) Delete(data int) *AVLNode {
	if node == nil {
		return nil
	}
	if data == node.data {
		if node.left == nil && node.right == nil {
			return nil
		}
		if node.left != nil && node.right != nil {
			if node.left.height > node.right.height {
				maxNode := node.left
				for maxNode.right != nil {
					maxNode = maxNode.right
				}
				node.data = maxNode.data
				node.left = node.left.Delete(maxNode.data)
				node.UpdateHeight()
			} else {
				minNode := node.right
				for minNode.left != nil {
					minNode = minNode.left
				}
				node.data = minNode.data
				node.right = node.right.Delete(minNode.data)
				node.UpdateHeight()
			}
		} else {
			if node.left != nil {
				node.data = node.left.data
				node.height = 1
				node.left = nil
			} else if node.right != nil {
				node.data = node.right.data
				node.height = 1
				node.right = nil
			}
		}
		return node
	} else if data < node.data {
		node.left = node.left.Delete(data)
		node.left.UpdateHeight()
	} else if data > node.data {
		node.right = node.right.Delete(data)
		node.right.UpdateHeight()
	}
	var newNode *AVLNode
	if node.BalanceFactor() == 2 {
		if node.left.BalanceFactor() >= 0 {
			newNode = node.RightRotate()
		} else {
			newNode = node.LeftRightNode()
		}
	} else if node.BalanceFactor() == -2 {
		if node.right.BalanceFactor() <= 0 {
			newNode = node.LeftRotate()
		} else {
			newNode = node.RightLeftNode()
		}
	}
	if newNode == nil {
		node.UpdateHeight()
		return node
	} else {
		//newNode.UpdateHeight()
		return newNode
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
