package RBtree

const (
	RED   = true
	BLACK = false
)

type Int int

func (x Int) Less(than Item) bool {
	return x < than.(Int)
}

type Item interface {
	Less(than Item) bool
}

type RBNode struct {
	Parent *RBNode
	Left   *RBNode
	Right  *RBNode
	color  bool
	Item
}

type RBTree struct {
	NIL   *RBNode
	root  *RBNode
	count int
}

func NewRBTree() *RBTree {
	node := &RBNode{nil, nil, nil, BLACK, nil}
	return &RBTree{
		NIL:   node,
		root:  node,
		count: 0,
	}
}

func (r *RBTree) LeftRotate(node *RBNode) {
	if node.Right == r.NIL {
		return
	}
	pivot := node.Right
	pivotL := pivot.Left
	node.Right = pivotL
	if pivotL != r.NIL {
		pivotL.Parent = node
	}

	pivot.Parent = node.Parent
	if node.Parent == r.NIL {
		r.root = pivot
	} else if node == node.Parent.Left {
		node.Parent.Left = pivot
	} else {
		node.Parent.Right = pivot
	}
	pivot.Left = node
	node.Parent = pivot
}

func (r *RBTree) RightRotate(node *RBNode) {
	if node.Left == r.NIL {
		return
	}
	pivot := node.Left
	pivotR := pivot.Right
	node.Left = pivotR
	if pivotR != r.NIL {
		pivotR.Parent = node
	}
	pivot.Parent = node.Parent
	if node.Parent == r.NIL {
		r.root = pivot
	} else if node == node.Parent.Left {
		node.Parent.Left = pivot
	} else {
		node.Parent.Right = pivot
	}
	pivot.Right = node
	node.Parent = pivot
}

func (r *RBTree) Insert(item Item) *RBNode {
	if item == nil {
		return nil
	}
	return r.insert(&RBNode{r.NIL, r.NIL, r.NIL, RED, item})
}

func (r *RBTree) insert(node *RBNode) *RBNode {
	x := r.root
	y := r.NIL
	//二分查找对应节点的值
	for x != r.NIL {
		y = x
		if node.Less(x.Item) {
			x = x.Left
		} else if x.Less(node.Item) {
			x = x.Right
		} else {
			//值已经存在，则直接返回
			return x
		}
	}
	//将节点的父节点设置为查找到的节点
	node.Parent = y
	if y == r.NIL {
		r.root = node
	} else if node.Less(y.Item) {
		y.Left = node
	} else {
		y.Right = node
	}
	r.count++
	r.InsertFixUp(node)
	return node
}

// case 1:插入节点父节点为黑色，直接插入
func (r *RBTree) InsertFixUp(node *RBNode) {
	for node.Parent.color == RED {
		//case 2：插入节点为父父节点的左侧
		if node.Parent == node.Parent.Parent.Left {
			//case 2.1：插入父节点的兄弟节点为红色，
			//将父节点和其兄弟节点变为黑色，父父节点变为红色，迭代更新
			y := node.Parent.Parent.Right
			if y.color == RED {
				node.Parent.color = BLACK
				y.color = BLACK
				node.Parent.Parent.color = RED
				node = node.Parent.Parent
			} else {
				//case 2.2：父节点的兄弟节点为黑色，且插入节点在父节点的右侧
				// 则进行以父节点左旋，转换为case2.3
				if node == node.Parent.Right {
					node = node.Parent
					r.LeftRotate(node)
				}
				// case 2.3：父节点的兄弟节点为黑色，且插入节点为父节点左侧
				// 则将父节点变为黑色，父父节点变为红色，以父父节点进行右旋
				node.Parent.color = BLACK
				node.Parent.Parent.color = RED
				r.RightRotate(node.Parent.Parent)
			}
		} else {
			y := node.Parent.Parent.Left
			if y.color == RED {
				node.Parent.color = BLACK
				y.color = BLACK
				node.Parent.Parent.color = RED
				node = node.Parent.Parent
			} else {
				if node == node.Parent.Left {
					node = node.Parent
					r.RightRotate(node)
				}
				node.Parent.color = BLACK
				node.Parent.Parent.color = RED
				r.LeftRotate(node.Parent.Parent)
			}
		}
	}
	r.root.color = BLACK
}

func (r *RBTree) min(node *RBNode) *RBNode {
	if node == r.NIL {
		return r.NIL
	}
	for node.Left != r.NIL {
		node = node.Left
	}
	return node
}

func (r *RBTree) max(node *RBNode) *RBNode {
	if node == r.NIL {
		return r.NIL
	}
	for node.Right != r.NIL {
		node = node.Right
	}
	return node
}

func (r *RBTree) Search(node *RBNode) *RBNode {
	p := r.root
	for p != r.NIL {
		if p.Less(node.Item) {
			p = p.Right
		} else if node.Less(p.Item) {
			p = p.Left
		} else {
			break
		}
	}
	return p
}

func (r *RBTree) successor(node *RBNode) *RBNode {
	if node == r.NIL {
		return r.NIL
	}

	if node.Right != r.NIL {
		return r.min(node.Right)
	}
	y := node.Parent
	for y != r.NIL && node == y.Right {
		node = y
		y = y.Parent
	}
	return y
}

func (r *RBTree) Delete(item Item) Item {
	if item == nil {
		return nil
	}
	return r.delete(&RBNode{
		Parent: r.NIL,
		Left:   r.NIL,
		Right:  r.NIL,
		color:  RED,
		Item:   item,
	}).Item
}

func (r *RBTree) delete(node *RBNode) *RBNode {
	//搜索节点位置
	z := r.Search(node)
	if z == r.NIL {
		return r.NIL
	}
	ret := &RBNode{
		Parent: r.NIL,
		Left:   r.NIL,
		Right:  r.NIL,
		color:  z.color,
		Item:   z.Item,
	}

	var x, y *RBNode
	//如果节点只有一个左或右节点，则使用子节点替换节点
	if z.Left == r.NIL || z.Right == r.NIL {
		y = z
	} else {
		//如果两个几点都存在，则使用右子树最小值替换
		y = r.successor(z)
	}

	//删除替换的节点
	if y.Left != r.NIL {
		x = y.Left
	} else {
		x = y.Right
	}

	x.Parent = y.Parent

	//如果y是根节点，则更新根节点
	if y.Parent == r.NIL {
		r.root = x
	} else if y == y.Parent.Left {
		//删除当前节点
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}

	//替换寻找到的节点值
	if y != z {
		z.Item = y.Item
	}

	//如果寻找到的替换的值是黑色节点，则需要调整颜色
	if y.color == BLACK {
		r.deleteFixUp(x)
	}
	r.count--
	return ret
}

// case 1：替换节点是红色节点，删除不影响红黑树平衡，直接替换即可
func (r *RBTree) deleteFixUp(node *RBNode) {
	for node != r.root && node.color == BLACK {
		//case 2：替换节点为黑色，且位于其父节点的左子节点
		if node == node.Parent.Left {
			//case 2.1：替换节点的兄弟节点为红色，设置父节点为红色，兄弟节点为黑色
			//再通过父节点左旋调整平衡
			w := node.Parent.Right
			if w.color == RED {
				w.color = BLACK
				node.Parent.color = RED
				r.LeftRotate(node.Parent)
				//迭代处理
				w = node.Parent.Right
			}
			// case 2.2：替换节点的兄弟节点为黑色，且其兄弟节点的左右节点均为黑色
			// 将兄弟节点设置为红色，以父节点作为替换节点进行处理
			if w.Left.color == BLACK && w.Right.color == BLACK {
				w.color = RED
				node = node.Parent
			} else {
				//case 2.3：兄弟节点为黑色，且兄弟节点的左节点为红色
				// 将兄弟节点的左节点设置为黑色，兄弟节点设置为红色，
				//以兄弟节点右旋，到达情况case 2.4
				if w.Right.color == BLACK {
					w.Left.color = BLACK
					w.color = RED
					r.RightRotate(w)
					w = node.Parent.Right
				}
				// case 2.4：兄弟节点为黑色，兄弟节点的右节点为黑色，
				// 将兄弟节点的右节点设置为黑色，以父节点进行左旋
				w.color = node.Parent.color
				node.Parent.color = BLACK
				w.Right.color = BLACK
				r.LeftRotate(node.Parent)
				node = r.root
			}
		} else {
			w := node.Parent.Left
			if w.color == RED {
				w.color = BLACK
				node.Parent.color = RED
				r.RightRotate(node.Parent)
				w = node.Parent.Left
			}
			if w.Left.color == BLACK && w.Right.color == BLACK {
				w.color = RED
				node = node.Parent
			} else {
				if w.Left.color == BLACK {
					w.Right.color = BLACK
					w.color = RED
					r.LeftRotate(w)
					w = node.Parent.Left
				}
				w.color = node.Parent.color
				node.Parent.color = BLACK
				w.Left.color = BLACK
				r.RightRotate(node.Parent)
				node = r.root
			}
		}
	}
	node.color = BLACK
}

func (r *RBTree) Len() int {
	return r.count
}
