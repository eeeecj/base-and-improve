package BPTree

type BPItem struct {
	key   int
	value interface{}
}
type BPNode struct {
	isLeaf   bool
	isRoot   bool
	parent   *BPNode
	pre      *BPNode
	next     *BPNode
	entries  []*BPItem
	children []*BPNode
}

func (node *BPNode) findItem(key int) interface{} {
	left, right := 0, len(node.entries)-1
	for left <= right {
		mid := left + (right-left)/2
		if node.entries[mid].key == key {
			return node.entries[mid].value
		} else if node.entries[mid].key > key {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return nil
}

func (node *BPNode) findNode(child *BPNode) int {
	for i := 0; i < len(node.children); i++ {
		if node.children[i] == child {
			return i
		}
	}
	return -1
}
func NewBPNode(order int, isLeaf, isRoot bool) *BPNode {
	if isLeaf {
		return &BPNode{
			isLeaf:  isLeaf,
			isRoot:  isRoot,
			entries: make([]*BPItem, 0, order+1),
		}
	} else {
		return &BPNode{
			isLeaf:   isLeaf,
			isRoot:   isRoot,
			children: make([]*BPNode, 0, order+1),
			entries:  make([]*BPItem, 0, order+1),
		}
	}
}

type BPTree struct {
	root   *BPNode
	order  int
	head   *BPNode
	height int
}

func NewBPTree(order int) *BPTree {
	if order < 3 {
		order = 3
	}
	return &BPTree{
		root:   NewBPNode(order, true, true),
		order:  order,
		head:   nil,
		height: 0,
	}
}

func (t *BPTree) Get(key int) interface{} {
	return t.root.get(key)
}

func (node *BPNode) get(key int) interface{} {
	if node.isLeaf {
		return node.findItem(key)
	}
	if node.entries[0].key > key {
		return node.children[0].get(key)
	} else if node.entries[len(node.entries)-1].key <= key {
		return node.children[len(node.children)-1].get(key)
	} else {
		left, right := 0, len(node.children)-1
		for left <= right {
			mid := left + (right-left)/2
			if node.entries[mid].key == key {
				return node.children[mid+1].get(key)
			} else if node.entries[mid].key > key {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		return node.children[left].get(key)
	}
}
func (t *BPTree) Insert(key int, value interface{}) {
	t.root.insert(key, value, t)
}
func (node *BPNode) insert(key int, value interface{}, tree *BPTree) {
	item := &BPItem{key: key, value: value}
	if node.isLeaf {
		node.addItem(item)
		if tree.height == 0 {
			tree.height = 1
		}
		if len(node.entries) < tree.order {
			return
		}
		node.splitNode(tree)
		return
	}
	if node.entries[0].key > key {
		node.children[0].insert(key, value, tree)
	} else if node.entries[len(node.entries)-1].key <= key {
		node.children[len(node.entries)-1].insert(key, value, tree)
	} else {
		left, right := 0, len(node.entries)-1
		for left <= right {
			mid := left + (right-left)/2
			if node.entries[mid].key == key {
				node.children[mid+1].insert(key, value, tree)
				return
			} else if node.entries[mid].key < key {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		node.children[left].insert(key, value, tree)
	}
}

func (node *BPNode) splitNode(tree *BPTree) {
	left := NewBPNode(tree.order, true, false)
	right := NewBPNode(tree.order, true, false)
	if node.pre != nil {
		node.pre.next = left
		left.pre = node.pre
	}
	if node.next != nil {
		node.next.pre = right
		right.next = node.next
	}
	if node.pre == nil {
		tree.head = left
	}
	left.next = right
	right.pre = left
	node.pre = nil
	node.next = nil
	node.splitLeafNode(left, right)

	if node.parent != nil {
		idx := node.parent.findNode(node)
		left.parent = node.parent
		right.parent = node.parent
		tempChild := append([]*BPNode{left, right}, node.parent.children[idx+1:]...)
		node.parent.children = append(node.parent.children[:idx], tempChild...)
		node.parent.entries = append(append(node.parent.entries[:idx], right.entries[0]), node.parent.entries[idx:]...)

		node.parent.splitInternalNode(tree)
		node.parent = nil
	} else {
		node.isRoot = false
		parent := NewBPNode(tree.order, false, true)
		tree.root = parent
		tree.height += 1
		left.parent = parent
		right.parent = parent
		parent.children = append(parent.children, left)
		parent.children = append(parent.children, right)
		parent.entries = append(parent.entries, right.entries[0])
	}
}

func (node *BPNode) splitInternalNode(tree *BPTree) {
	m := len(node.children)
	if m > tree.order {
		left := NewBPNode(tree.order, false, false)
		right := NewBPNode(tree.order, false, false)
		for i := 0; i < m; i++ {
			if i < (m+1)/2 {
				left.children = append(left.children, node.children[i])
				node.children[i].parent = left
				left.entries = append(left.entries, node.entries[i])
			} else {
				right.children = append(right.children, node.children[i])
				node.children[i] = right
				right.entries = append(right.entries, node.entries[i])
			}
		}
		if node.parent != nil {
			idx := node.parent.findNode(node)
			left.parent = node.parent
			right.parent = node.parent
			tempChild := append([]*BPNode{left, right}, node.parent.children[idx+1:]...)
			node.parent.children = append(node.parent.children[:idx], tempChild...)
			node.parent.entries = append(append(node.parent.entries[:idx], right.entries[0]), node.parent.entries[idx:]...)

			node.parent.splitInternalNode(tree)
			node.parent = nil
		} else {
			node.isRoot = false
			parent := NewBPNode(tree.order, false, true)
			tree.root = parent
			tree.height += 1
			right.parent = parent
			left.parent = parent
			parent.children = append(parent.children, left)
			parent.children = append(parent.children, right)
			parent.entries = append(parent.entries, right.entries[0])
		}
	}
}

func (node *BPNode) splitLeafNode(left *BPNode, right *BPNode) {
	for i := 0; i < len(node.entries); i++ {
		if i <= (len(node.entries))/2 {
			left.entries = append(left.entries, node.entries[i])
		} else {
			right.entries = append(right.entries, node.entries[i])
		}
	}
}

func (t *BPTree) Delete(key int) interface{} {
	return t.root.delete(key, t)
}

func (node *BPNode) delete(key int, tree *BPTree) interface{} {
	if node.isLeaf {
		if node.findItem(key) == nil {
			return nil
		}
		if node.isRoot {
			if len(node.entries) == 0 {
				tree.height = 0
			}
			return node.deleteItem(key)
		}
		if len(node.entries) > tree.order/2 && len(node.entries) > 2 {
			return node.deleteItem(key)
		}
		if node.pre != nil && node.pre.parent == node.parent &&
			len(node.pre.entries) > tree.order/2 && len(node.pre.entries) > 2 {
			m := len(node.pre.entries)
			node.entries = append([]*BPItem{node.pre.entries[m-1]}, node.entries...)
			node.pre.entries = node.pre.entries[:m-1]
			idx := node.parent.findNode(node.pre)
			node.parent.entries[idx] = node.entries[0]
			return node.deleteItem(key)
		}

		if node.next != nil && node.next.parent == node.parent &&
			len(node.next.entries) > tree.order/2 && len(node.next.entries) > 2 {
			node.entries = append(node.entries, node.next.entries[0])
			node.next.entries = node.next.entries[1:]
			idx := node.parent.findNode(node.next)
			node.parent.entries[idx] = node.next.entries[0]
			return node.deleteItem(key)
		}
		if node.pre != nil && node.pre.parent == node.parent &&
			(len(node.pre.entries) <= tree.order/2 || len(node.pre.entries) <= 2) {
			res := node.deleteItem(key)
			node.mergePreNode(node.pre, tree)
			if (!node.parent.isRoot && len(node.parent.children) >= tree.order/2 &&
				len(node.parent.children) >= 2) || (node.parent.isRoot && len(node.parent.children) >= 2) {
				return res
			}
			node.parent.updateRemove(tree)
			return res
		}

		if node.next != nil && node.next.parent == node.parent &&
			(len(node.next.entries) <= tree.order/2 || len(node.next.entries) <= 2) {
			res := node.deleteItem(key)
			node.mergeNexNode(node.next, tree)
			if (!node.parent.isRoot && (len(node.parent.children) >= tree.order/2 &&
				len(node.parent.children) >= 2)) || (node.parent.isRoot && len(node.parent.children) >= 2) {
				return res
			}
			node.parent.updateRemove(tree)
			return res
		}
	}
	if node.entries[0].key > key {
		return node.children[0].delete(key, tree)
	} else if node.entries[len(node.entries)-1].key <= key {
		return node.children[len(node.children)-1].delete(key, tree)
	} else {
		left, right := 0, len(node.entries)-1
		for left <= right {
			mid := left + (right-left)/2
			if node.entries[mid].key == key {
				return node.children[mid+1].delete(key, tree)
			} else if node.entries[mid].key > key {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		return node.children[left].delete(key, tree)
	}
}
func (node *BPNode) mergeNexNode(next *BPNode, tree *BPTree) {
	for i := 0; i < len(next.entries); i++ {
		node.entries = append(node.entries, next.entries[i])
	}
	nidx := node.parent.findNode(next)
	node.parent.children = append(node.parent.children[:nidx], node.parent.children[nidx+1:]...)
	if next.next != nil {
		temp := next
		temp.next.pre = node
		node.next = temp.next
		temp.pre = nil
		temp.next = nil
	} else {
		next.pre = nil
		node.next = nil
	}
	idx := node.parent.findNode(node)
	node.parent.entries = append(node.parent.entries[:idx], node.parent.entries[idx+1:]...)
}

func (node *BPNode) mergePreNode(pre *BPNode, tree *BPTree) {
	for i := 0; i < len(node.entries); i++ {
		pre.entries = append(pre.entries, node.entries[i])
	}
	node.entries = pre.entries
	idx := node.parent.findNode(pre)
	node.parent.children = append(node.parent.children[:idx], node.parent.children[idx+1:]...)
	pre.parent = nil

	if pre.pre != nil {
		temp := pre
		temp.pre.next = node
		node.pre = temp.pre
		temp.pre = nil
		temp.next = nil
	} else {
		tree.head = node
		pre.next = nil
		node.pre = nil
	}
	idx = node.parent.findNode(node)
	node.parent.entries = append(node.parent.entries[:idx], node.parent.entries[idx+1:]...)
}

func (node *BPNode) updateRemove(tree *BPTree) {
	if len(node.children) < tree.order/2 || len(node.children) < 2 {
		if node.isRoot {
			if len(node.children) >= 2 {
				return
			}
			root := node.children[0]
			tree.root = root
			tree.height -= 1
			return
		}
		curIdx := node.parent.findNode(node)
		preIdx := curIdx - 1
		nextIdx := curIdx + 1
		var pre, next *BPNode
		if preIdx > 0 {
			pre = node.parent.children[preIdx]
		}
		if nextIdx < len(node.parent.children) {
			next = node.parent.children[nextIdx]
		}

		if pre != nil && len(pre.children) > tree.order/2 && len(pre.children) > 2 {
			m := len(pre.children)
			bw := pre.children[m-1]
			pre.children = pre.children[:m-1]
			bw.parent = node
			node.children = append([]*BPNode{bw}, node.children...)

			pidx := node.parent.findNode(pre)
			node.entries = append([]*BPItem{node.parent.entries[pidx]}, node.entries...)
			node.parent.entries[pidx] = pre.entries[m-2]
			pre.entries = pre.entries[:m-2]
			return
		}

		if next != nil && len(next.children) > tree.order/2 && len(next.children) > 2 {
			bw := next.children[0]
			next.children = next.children[1:]
			bw.parent = node
			node.children = append(node.children, bw)

			pidx := node.parent.findNode(node)
			node.entries = append(node.entries, node.parent.entries[pidx])
			node.parent.entries[pidx] = next.entries[0]
			pre.entries = pre.entries[1:]
			return
		}

		if pre != nil && len(pre.children) <= tree.order/2 || len(pre.children) <= 2 {
			for i := 0; i < len(node.children); i++ {
				pre.children = append(pre.children, node.children[i])
			}

			for i := 0; i < len(pre.children); i++ {
				pre.children[i].parent = node
			}
			idx := node.parent.findNode(pre)

			pre.entries = append(pre.entries, node.parent.entries[idx])
			for i := 0; i < len(node.entries); i++ {
				pre.entries = append(pre.entries, node.entries[i])
			}

			node.children = pre.children
			node.entries = pre.entries

			node.parent.children = append(node.parent.children[:idx], node.parent.children[idx+1:]...)
			nidx := node.parent.findNode(node)
			node.parent.entries = append(node.parent.entries[:nidx], node.parent.entries[nidx+1:]...)

			if (!node.parent.isRoot && (len(node.parent.children) >= tree.order/2 &&
				len(node.parent.children) >= 2)) || (node.parent.isRoot && len(node.parent.children) >= 2) {
				return
			}
			node.parent.updateRemove(tree)
			return
		}

		if next != nil && (len(next.children) <= tree.order/2 || len(next.children) <= 2) {
			for i := 0; i < len(next.children); i++ {
				child := next.children[i]
				node.children = append(node.children, child)
				child.parent = node
			}
			idx := node.parent.findNode(node)
			node.entries = append(node.entries, node.parent.entries[idx])
			for i := 0; i < len(next.entries); i++ {
				node.entries = append(node.entries, next.entries[i])
			}
			nidx := node.parent.findNode(next)
			node.parent.children = append(node.parent.children[:nidx], node.parent.children[nidx+1:]...)
			node.parent.entries = append(node.parent.entries[:idx], node.parent.entries[idx+1:]...)

			if (!node.parent.isRoot && (len(node.parent.children) >= tree.order/2 && len(node.parent.children) >= 2)) ||
				(node.parent.isRoot && len(node.parent.children) >= 2) {
				return
			}
			node.parent.updateRemove(tree)
		}
	}
}

func (node *BPNode) addItem(item *BPItem) {
	left, right := 0, len(node.entries)-1
	for left <= right {
		mid := left + (right-left)/2
		if node.entries[mid].key == item.key {
			node.entries[mid] = item
			return
		} else if node.entries[mid].key > item.key {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	node.entries = append(node.entries, nil)
	copy(node.entries[left+1:], node.entries[left:])
	node.entries[left] = item
}

func (node *BPNode) deleteItem(key int) interface{} {
	left, right := 0, len(node.entries)-1
	for left <= right {
		mid := left + (right-left)/2
		if node.entries[mid].key == key {
			v := node.entries[mid]
			node.entries = append(node.entries[:mid], node.entries[mid+1:]...)
			return v
		} else if node.entries[mid].key > key {
			right = mid + 1
		} else {
			left = mid - 1
		}
	}
	return nil
}
