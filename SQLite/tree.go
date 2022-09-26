package main

import (
	"fmt"
	"unsafe"
)

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

const (
	LEAF_NODE_NUM_CELL_SIZE   = 4
	LEAF_NODE_NUM_CELL_OFFSET = COMMON_NODE_HEADER_SIZE
	LEAF_NODE_HEADER_SIZE     = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELL_SIZE
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

// 此处指针操作需要注意
func leafNodeNumCells(node unsafe.Pointer) (int, unsafe.Pointer) {
	p := (*[LEAF_NODE_HEADER_SIZE]byte)(node)
	cellNums := make([]byte, LEAF_NODE_CELL_SIZE)
	copy(cellNums[:LEAF_NODE_CELL_SIZE], p[LEAF_NODE_NUM_CELL_OFFSET:LEAF_NODE_HEADER_SIZE])
	cell_ids := int(BytesToInt32(cellNums))
	return cell_ids, unsafe.Pointer(&p[0])
}
func leafNodeNumCellsSet(node unsafe.Pointer, num int) {
	p := (*[LEAF_NODE_HEADER_SIZE]byte)(node)
	copy(p[LEAF_NODE_NUM_CELL_OFFSET:LEAF_NODE_HEADER_SIZE], Uint32ToBytes(int32(num)))
}

func leafNodeCell(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	p := (*[LEAF_NODE_SPACE_FOR_CELLS]byte)(node)
	return unsafe.Pointer(&p[LEAF_NODE_HEADER_SIZE+cellNum*LEAF_NODE_CELL_SIZE])
}

func leafNodeKey(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	pp := leafNodeCell(node, cellNum)
	return pp
}
func leafNodeKeySet(node unsafe.Pointer, num int) {
	p := (*[LEAF_NODE_KEY_SIZE]byte)(node)
	copy(p[0:LEAF_NODE_KEY_SIZE], Uint32ToBytes(int32(num)))
}

func leafNodeValue(node unsafe.Pointer, cellNum int) unsafe.Pointer {
	pp := leafNodeCell(node, cellNum)
	p := (*[LEAF_NODE_CELL_SIZE]byte)(pp)
	return unsafe.Pointer(&p[LEAF_NODE_KEY_SIZE])
}

func initLeafNode(node unsafe.Pointer) {
	_, pp := leafNodeNumCells(node)
	leafNodeNumCellsSet(pp, 0)
}

func leafNodeInsert(cursor *Cursor, key int, value *Row) {
	node := getPage(cursor.table.Pager, cursor.pageNum)
	numCells, _ := leafNodeNumCells(node)
	if numCells >= LEAF_NODE_MAX_CELLS {
		fmt.Println("Need to implement splitting a leaf node.")
		return
	}
	if cursor.cellNum < numCells {
		for i := numCells; i > cursor.cellNum; i-- {
			dest := *(*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(node, i))
			origin := *(*[LEAF_NODE_CELL_SIZE]byte)(leafNodeCell(node, i-1))
			copy(dest[0:LEAF_NODE_CELL_SIZE], origin[0:LEAF_NODE_CELL_SIZE])
		}
	}
	cells, p := leafNodeNumCells(node)
	leafNodeNumCellsSet(p, cells+1)
	leafNodeKeySet(leafNodeKey(node, cursor.cellNum), key)
	serializeRow(value, leafNodeValue(node, cursor.cellNum))
}
