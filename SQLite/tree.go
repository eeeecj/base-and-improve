package main

import (
	"fmt"
	"os"
	"unsafe"
)

//表格结构请查看链接：https://cstack.github.io/db_tutorial/assets/images/leaf-node-format.png

type NodeType int32

const (
	NODE_INTERNAL NodeType = iota
	NODE_LEAF
)

const (
	NODE_TYPE_SIZE          = 1
	NODE_TYPE_OFFSET        = 0
	IS_ROOT_SIZE            = 1
	IS_ROOT_OFFSET          = NODE_TYPE_SIZE
	PARENT_POINTER_SIZE     = 4
	PARENT_POINTER_OFFSET   = IS_ROOT_OFFSET + IS_ROOT_SIZE
	COMMON_NODE_HEADER_SIZE = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE
)

// 叶子节点格式：https://cstack.github.io/db_tutorial/assets/images/leaf-node-format.png
const (
	LEAF_NODE_NUM_CELL_SIZE   = 4
	LEAF_NODE_NUM_CELL_OFFSET = COMMON_NODE_HEADER_SIZE
	//LEAF_NODE_HEADER_SIZE     = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELL_SIZE
	LEAF_NODE_NEXT_LEAF_SIZE   = 4
	LEAF_NODE_NEXT_LEAF_OFFSET = LEAF_NODE_NUM_CELL_OFFSET + LEAF_NODE_NUM_CELL_SIZE
	LEAF_NODE_HEADER_SIZE      = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELL_SIZE + LEAF_NODE_NEXT_LEAF_SIZE
)

const (
	LEAF_NODE_KEY_SIZE        = 4
	LEAF_NODE_KEY_OFFSET      = 0
	LEAF_NODE_VALUE_SIZE      = ROW_SIZE
	LEAF_NODE_VALUE_OFFSET    = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE
	LEAF_NODE_CELL_SIZE       = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE
	LEAF_NODE_SPACE_FOR_CELLS = PAGE_SIZE - LEAF_NODE_HEADER_SIZE
	LEAF_NODE_MAX_CELLS       = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE
)

const (
	LEAF_NODE_RIGHT_SPLIT_COUNT = (LEAF_NODE_MAX_CELLS + 1) / 2
	LEAF_NODE_LEFT_SPLIT_COUNT  = (LEAF_NODE_MAX_CELLS + 1) - LEAF_NODE_RIGHT_SPLIT_COUNT
)

// 内部节点格式：https://cstack.github.io/db_tutorial/assets/images/internal-node-format.png
const (
	INTERNAL_NODE_MAX_CELLS          = 3
	INTERNAL_NODE_NUM_KEYS_SIZE      = 4
	INTERNAL_NODE_NUM_KEYS_OFFSET    = COMMON_NODE_HEADER_SIZE
	INTERNAL_NODE_RIGHT_CHILD_SIZE   = 4
	INTERNAL_NODE_RIGHT_CHILD_OFFSET = INTERNAL_NODE_NUM_KEYS_OFFSET + INTERNAL_NODE_NUM_KEYS_SIZE
	INTERNAL_NODE_HEADER_SIZE        = COMMON_NODE_HEADER_SIZE + INTERNAL_NODE_NUM_KEYS_SIZE + INTERNAL_NODE_RIGHT_CHILD_SIZE
)

const (
	INTERNAL_NODE_KEY_SIZE   = 4
	INTERNAL_NODE_CHILD_SIZE = 4
	INTERNAL_NODE_CELL_SIZE  = INTERNAL_NODE_CHILD_SIZE + INTERNAL_NODE_KEY_SIZE
)

func getNodeType(node unsafe.Pointer) NodeType {
	p := (*[NODE_TYPE_SIZE]byte)(unsafe.Pointer(uintptr(node) + NODE_TYPE_OFFSET))
	return NodeType(p[0])
}

func setNodeType(node unsafe.Pointer, t NodeType) {
	p := (*[NODE_TYPE_SIZE]byte)(unsafe.Pointer(uintptr(node) + NODE_TYPE_OFFSET))
	ty := Uint32ToBytes(int32(t))
	copy(p[:], ty[len(ty)-1:])
}

// 判断是否根节点
func isNodeRoot(node unsafe.Pointer) bool {
	pp := unsafe.Pointer(uintptr(node) + IS_ROOT_OFFSET)
	value := *(*bool)(pp)
	return value
}

func setNodeRoot(node unsafe.Pointer, isRoot bool) {
	if isRoot {
		*(*uint8)(unsafe.Pointer(uintptr(node) + IS_ROOT_OFFSET)) = 1
	} else {
		*(*uint8)(unsafe.Pointer(uintptr(node) + IS_ROOT_OFFSET)) = 0
	}
}

// 设置父节点
func nodeParent(node unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(node) + PARENT_POINTER_OFFSET)
}
func nodeParentSet(node unsafe.Pointer, num int) {
	p := unsafe.Pointer(uintptr(node) + PARENT_POINTER_OFFSET)
	copy((*[PARENT_POINTER_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取内部节点存储的键值对数量
func internalNodeNumKeys(node unsafe.Pointer) (int, unsafe.Pointer) {
	p := unsafe.Pointer(uintptr(node) + INTERNAL_NODE_NUM_KEYS_OFFSET)
	res := BytesToInt32((*[INTERNAL_NODE_NUM_KEYS_SIZE]byte)(p)[:])
	return int(res), p
}

func internalNodeNumKeysSet(node unsafe.Pointer, num int) {
	p := unsafe.Pointer(uintptr(node) + INTERNAL_NODE_NUM_KEYS_OFFSET)
	copy((*[INTERNAL_NODE_NUM_KEYS_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取内部节点的右子树
func internalNodeRightChild(node unsafe.Pointer) (int, unsafe.Pointer) {
	p := unsafe.Pointer(uintptr(node) + INTERNAL_NODE_RIGHT_CHILD_OFFSET)
	res := BytesToInt32((*[INTERNAL_NODE_RIGHT_CHILD_SIZE]byte)(p)[:])
	return int(res), p
}

// 获取对应内部节点的节点值
func internalNodeCell(node unsafe.Pointer, cellNum int) (int, unsafe.Pointer) {
	p := unsafe.Pointer(uintptr(node) + INTERNAL_NODE_HEADER_SIZE + uintptr(cellNum*INTERNAL_NODE_CELL_SIZE))
	res := BytesToInt32((*[INTERNAL_NODE_KEY_SIZE]byte)(p)[:])
	return int(res), p
}

// 获取对应内部节点的子节点
func internalNodeChild(node unsafe.Pointer, childNum int) (int, unsafe.Pointer) {
	numKey, _ := internalNodeNumKeys(node)
	if childNum > numKey {
		fmt.Printf("Tried to access child_num %d > num_keys %d\n", childNum, numKey)
		os.Exit(0)
		return 0, nil
	} else if childNum == numKey {
		//当查找的子节点大于当前节点最大的值，则从右节点查找
		return internalNodeRightChild(node)
	} else {
		return internalNodeCell(node, childNum)
	}
}

// 获取内部节点key
func internalNodeKey(node unsafe.Pointer, keyNum int) (int, unsafe.Pointer) {
	_, pp := internalNodeCell(node, keyNum)
	p := unsafe.Pointer(uintptr(pp) + INTERNAL_NODE_CHILD_SIZE)
	res := BytesToInt32((*[INTERNAL_NODE_KEY_SIZE]byte)(p)[:])
	return int(res), p
}

func internalNodeKeySet(node unsafe.Pointer, keyNum int, num int) {
	_, pp := internalNodeCell(node, keyNum)
	p := unsafe.Pointer(uintptr(pp) + INTERNAL_NODE_CHILD_SIZE)
	copy((*[INTERNAL_NODE_KEY_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取叶节点的数量
func leafNodeNumCells(node unsafe.Pointer) (int, unsafe.Pointer) {
	pp := unsafe.Pointer(uintptr(node) + LEAF_NODE_NUM_CELL_OFFSET)
	res := BytesToInt32((*[LEAF_NODE_NUM_CELL_SIZE]byte)(pp)[:])
	return int(res), pp
}

func leafNodeNumCellsSet(node unsafe.Pointer, num int) {
	p := unsafe.Pointer(uintptr(node) + LEAF_NODE_NUM_CELL_OFFSET)
	copy((*[LEAF_NODE_NUM_CELL_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取叶节点的叶节点
func leafNodeNextLeaf(node unsafe.Pointer) (int, unsafe.Pointer) {
	pp := unsafe.Pointer(uintptr(node) + LEAF_NODE_NEXT_LEAF_OFFSET)
	res := BytesToInt32((*[LEAF_NODE_NEXT_LEAF_SIZE]byte)(pp)[:])
	return int(res), pp
}

func leafNodeNextLeafSet(node unsafe.Pointer, num int) {
	p := unsafe.Pointer(uintptr(node) + LEAF_NODE_NEXT_LEAF_OFFSET)
	copy((*[LEAF_NODE_NUM_CELL_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取叶节点对应节点的节点值
func leafNodeCell(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	p := unsafe.Pointer(uintptr(node) + LEAF_NODE_HEADER_SIZE + uintptr(cellNum*LEAF_NODE_CELL_SIZE))
	return p
}

// 获取叶节点对应节点的key
func leafNodeKey(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	pp := leafNodeCell(node, cellNum)
	return pp
}
func leafNodeKeySet(node unsafe.Pointer, cellNum int, num int) {
	p := leafNodeCell(node, cellNum)
	copy((*[LEAF_NODE_KEY_SIZE]byte)(p)[:], Uint32ToBytes(int32(num)))
}

// 获取叶节点对应节点的值
func leafNodeValue(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	pp := unsafe.Pointer(uintptr(leafNodeCell(node, cellNum)) + LEAF_NODE_KEY_SIZE)
	return pp
}

// 获取节点对应的最大节点的key
func getNodeMaxKey(node unsafe.Pointer) int {
	switch getNodeType(node) {
	case NODE_INTERNAL:
		keys, _ := internalNodeNumKeys(node)
		res, _ := internalNodeKey(node, keys-1)
		return res
	case NODE_LEAF:
		keys, _ := leafNodeNumCells(node)
		res := BytesToInt32((*[LEAF_NODE_KEY_SIZE]byte)(leafNodeKey(node, keys-1))[:])
		return int(res)
	}
	return 0
}

// 初始化叶子节点
func initLeafNode(node unsafe.Pointer) {
	//设置叶子节点类型和
	setNodeType(node, NODE_LEAF)
	setNodeRoot(node, false)
	leafNodeNumCellsSet(node, 0)
	leafNodeNextLeafSet(node, 0)
}

func initInternalNode(node unsafe.Pointer) {
	setNodeType(node, NODE_INTERNAL)
	setNodeRoot(node, false)
	internalNodeNumKeysSet(node, 0)
}

// 从叶子节点查找对应键
func leafNodeFind(table *Table, pageNum int, key int) *Cursor {
	//获取页面
	node := getPage(table.Pager, pageNum)
	numCells, _ := leafNodeNumCells(node)
	cursor := &Cursor{}
	cursor.table = table
	cursor.pageNum = pageNum
	cursor.tableEnd = false

	//二分查找
	left, right := 0, numCells-1
	for left <= right {
		mid := left + (right-left)/2
		p := (*[LEAF_NODE_KEY_SIZE]byte)(leafNodeKey(node, mid))
		midKey := int(BytesToInt32(p[:]))
		if key == midKey {
			cursor.cellNum = mid
			return cursor
		} else if key >= midKey {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	cursor.cellNum = left
	return cursor
}

// 从内部节点查找子节点
func internalNodeFindChild(node unsafe.Pointer, key int) int {
	numKeys, _ := internalNodeNumKeys(node)
	//二分查找
	left, right := 0, numKeys-1
	for left <= right {
		mid := left + (right-left)/2
		keyRight, _ := internalNodeKey(node, mid)
		if keyRight <= key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

// 从内部节点查找
func internalNodeFind(table *Table, pageNum int, key int) *Cursor {
	node := getPage(table.Pager, pageNum)
	//根据键值大小查找对应分支
	childIndex := internalNodeFindChild(node, key)
	//从对应节点中查找对应的值
	childNum, _ := internalNodeChild(node, childIndex)
	//查找对应页面
	child := getPage(table.Pager, childNum)
	switch getNodeType(child) {
	case NODE_LEAF:
		//子节点则返回对应的值
		return leafNodeFind(table, childNum, key)
	case NODE_INTERNAL:
		//如果是中间节点继续二分查找
		return internalNodeFind(table, childNum, key)
	}
	return nil
}

// 返回键值存储的位置，如果不存在则返回应该插入的位置
func tableFind(table *Table, key int) *Cursor {
	rootPageNum := table.rootPageNum
	rootNode := getPage(table.Pager, rootPageNum)

	rootNodeType := getNodeType(rootNode)
	if rootNodeType == NODE_LEAF {
		return leafNodeFind(table, rootPageNum, key)
	} else {
		return internalNodeFind(table, rootPageNum, key)
	}
}

// 返回0对应的值所在的位置
func tableStart(table *Table) *Cursor {
	cursor := tableFind(table, 0)
	node := getPage(table.Pager, cursor.cellNum)
	numCells, _ := leafNodeNumCells(node)
	cursor.tableEnd = numCells == 0
	return cursor
}

// 获取未使用的页面
func getUnusedPageNum(pager *Pager) int {
	return pager.numPages
}

// 创建新的根节点，处理节点分裂，旧节点复制到新页面，变成左树，
// 新的树指向两个子树
func createNewRoot(table *Table, rightChildPageNum int) {
	//获取根几点
	root := getPage(table.Pager, table.rootPageNum)
	//获取右节点页面
	rightChild := getPage(table.Pager, rightChildPageNum)
	//获取左节点页面
	leftChildPageNum := getUnusedPageNum(table.Pager)
	leftChild := getPage(table.Pager, leftChildPageNum)
	//将根节点数据复制到左节点
	copy((*[PAGE_SIZE]byte)(leftChild)[:], (*[PAGE_SIZE]byte)(root)[:])
	setNodeRoot(leftChild, false)
	//初始化根节点
	initInternalNode(root)
	setNodeRoot(root, true)
	//设置根节点键数量
	internalNodeNumKeysSet(root, 1)
	//将根节点指向左节点
	_, pp := internalNodeChild(root, 0)
	copy((*[INTERNAL_NODE_CHILD_SIZE]byte)(pp)[:], Uint32ToBytes(int32(leftChildPageNum)))
	// 将左子树的最大值设置在根节点中方便查询
	leftChildMaxKey := getNodeMaxKey(leftChild)
	internalNodeKeySet(root, 0, leftChildMaxKey)
	// 将右节点对应的指针设置在根节点
	_, pp = internalNodeRightChild(root)
	copy((*[INTERNAL_NODE_RIGHT_CHILD_SIZE]byte)(pp)[:], Uint32ToBytes(int32(rightChildPageNum)))

	// 设置左节点和右节点的父节点为根节点
	nodeParentSet(nodeParent(leftChild), table.rootPageNum)
	nodeParentSet(nodeParent(rightChild), table.rootPageNum)
}

// 向对应的子节点的父节点插入键值对
func internalNodeInsert(table *Table, parentPageNum int, childPageNum int) {
	parent := getPage(table.Pager, parentPageNum)
	child := getPage(table.Pager, childPageNum)
	//从父结点中查找子节点位置
	childMaxKey := getNodeMaxKey(child)
	index := internalNodeFindChild(parent, childMaxKey)
	originalNumKey, _ := internalNodeNumKeys(parent)

	internalNodeNumKeysSet(parent, originalNumKey+1)

	if originalNumKey >= INTERNAL_NODE_MAX_CELLS {
		fmt.Printf("Need to implement splitting internal node\n")
		os.Exit(0)
	}

	// 获取右节点页面号
	rightChildPageNum, _ := internalNodeRightChild(parent)
	rightChild := getPage(table.Pager, rightChildPageNum)

	//如果子节点的键值大于右节点的最大值，则
	if childMaxKey > getNodeMaxKey(rightChild) {
		//将旧右节点存储在父节点中
		_, p := internalNodeChild(parent, originalNumKey)
		copy((*[INTERNAL_NODE_CHILD_SIZE]byte)(p)[:], Uint32ToBytes(int32(rightChildPageNum)))

		internalNodeKeySet(parent, originalNumKey, getNodeMaxKey(rightChild))

		//将右节点设置为较大的子节点值
		_, p = internalNodeRightChild(parent)
		copy((*[INTERNAL_NODE_RIGHT_CHILD_SIZE]byte)(p)[:], Uint32ToBytes(int32(childPageNum)))
	} else {
		//将节点右移保证大小顺序
		for i := originalNumKey; i > index; i-- {
			_, pp := internalNodeCell(parent, i)
			dest := (*[INTERNAL_NODE_CELL_SIZE]byte)(pp)
			_, pp = internalNodeCell(parent, i-1)
			source := (*[INTERNAL_NODE_CELL_SIZE]byte)(pp)
			copy(dest[:], source[:])
		}
		_, pp := internalNodeChild(parent, index)
		copy((*[INTERNAL_NODE_CHILD_SIZE]byte)(pp)[:], Uint32ToBytes(int32(childPageNum)))
		_, pp = internalNodeKey(parent, index)
		copy((*[INTERNAL_NODE_KEY_SIZE]byte)(pp)[:], Uint32ToBytes(int32(childMaxKey)))
	}
}

func updateInternalNodeKey(node unsafe.Pointer, oldKey int, newKey int) {
	oldChildIndex := internalNodeFindChild(node, oldKey)
	internalNodeKeySet(node, oldChildIndex, newKey)
}

// 创建新节点，移动一半节点，在两个节点间插入节点，更新或创建一个新父节点
func leafNodeSplitAndInsert(cursor *Cursor, key int, value *Row) {
	oldNode := getPage(cursor.table.Pager, cursor.pageNum)
	oldMax := getNodeMaxKey(oldNode)
	//创建新页面
	newPageNum := getUnusedPageNum(cursor.table.Pager)
	newNode := getPage(cursor.table.Pager, newPageNum)
	initLeafNode(newNode)
	nodeParentSet(newNode, int(BytesToInt32((*[PARENT_POINTER_SIZE]byte)(nodeParent(oldNode))[:])))
	//在旧页面与新页面中插入新页面
	oldNextPageNum, _ := leafNodeNextLeaf(oldNode)
	leafNodeNextLeafSet(newNode, oldNextPageNum)
	leafNodeNextLeafSet(oldNode, newPageNum)

	for i := LEAF_NODE_MAX_CELLS; i >= 0; i-- {
		var dest unsafe.Pointer
		if i >= LEAF_NODE_LEFT_SPLIT_COUNT {
			dest = newNode
		} else {
			dest = oldNode
		}
		indexWithinNode := i % LEAF_NODE_LEFT_SPLIT_COUNT
		dest = leafNodeCell(dest, indexWithinNode)

		if i == cursor.cellNum {
			//serializeRow(value, dest)
			serializeRow(value, leafNodeValue(dest, indexWithinNode))
			leafNodeKeySet(dest, indexWithinNode, key)
		} else if i > cursor.cellNum {
			d := (*[LEAF_NODE_CELL_SIZE]byte)(dest)
			o := (*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(oldNode, i-1))
			copy(d[:], o[:])
		} else {
			d := (*[LEAF_NODE_CELL_SIZE]byte)(dest)
			o := (*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(newNode, i-1))
			copy(d[:], o[:])
		}
	}
	leafNodeNumCellsSet(oldNode, LEAF_NODE_LEFT_SPLIT_COUNT)
	leafNodeNumCellsSet(newNode, LEAF_NODE_RIGHT_SPLIT_COUNT)

	if isNodeRoot(oldNode) {
		createNewRoot(cursor.table, newPageNum)
	} else {
		parentPageNum := int(BytesToInt32((*[PARENT_POINTER_SIZE]byte)(nodeParent(oldNode))[:]))
		newMax := getNodeMaxKey(oldNode)
		parent := getPage(cursor.table.Pager, parentPageNum)
		updateInternalNodeKey(parent, oldMax, newMax)
		internalNodeInsert(cursor.table, parentPageNum, newPageNum)
	}
}

func leafNodeInsert(cursor *Cursor, key int, value *Row) {
	node := getPage(cursor.table.Pager, cursor.pageNum)
	numCells, _ := leafNodeNumCells(node)
	if numCells >= LEAF_NODE_MAX_CELLS {
		//fmt.Println("Need to implement splitting a leaf node.")
		leafNodeSplitAndInsert(cursor, key, value)
		return
	}
	if cursor.cellNum < numCells {
		for i := numCells; i > cursor.cellNum; i-- {
			dest := (*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(node, i))
			origin := (*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(node, i-1))
			copy(dest[:], origin[:])
		}
	}
	cells, _ := leafNodeNumCells(node)
	leafNodeNumCellsSet(node, cells+1)
	leafNodeKeySet(node, cursor.cellNum, key)
	serializeRow(value, leafNodeValue(node, cursor.cellNum))
}
