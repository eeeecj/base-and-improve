package RBtree

import (
	"log"
	"testing"
)

func TestRBTree_LeftRotate(t *testing.T) {
	var i10 Int = 10
	var i12 Int = 12

	rbtree := NewRBTree()
	x := &RBNode{rbtree.NIL, rbtree.NIL, rbtree.NIL, BLACK, i10}
	rbtree.root = x
	y := &RBNode{rbtree.root.Right, rbtree.NIL, rbtree.NIL, RED, i12}
	rbtree.root.Right = y

	log.Println("root : ", rbtree.root)
	log.Println("left : ", rbtree.root.Left)
	log.Println("right : ", rbtree.root.Right)

	rbtree.LeftRotate(rbtree.root)

	log.Println("root : ", rbtree.root)
	log.Println("left : ", rbtree.root.Left)
	log.Println("right : ", rbtree.root.Right)
}

func TestRBTree_RightRotate(t *testing.T) {
	var i10 Int = 10
	var i12 Int = 12

	rbtree := NewRBTree()
	x := &RBNode{rbtree.NIL, rbtree.NIL, rbtree.NIL, BLACK, i10}
	rbtree.root = x
	y := &RBNode{rbtree.root.Right, rbtree.NIL, rbtree.NIL, RED, i12}
	rbtree.root.Left = y

	log.Println("root : ", rbtree.root)
	log.Println("left : ", rbtree.root.Left)
	log.Println("right : ", rbtree.root.Right)

	rbtree.RightRotate(rbtree.root)

	log.Println("root : ", rbtree.root)
	log.Println("left : ", rbtree.root.Left)
	log.Println("right : ", rbtree.root.Right)
}

func TestInsertAndDelete(t *testing.T) {
	rbt := NewRBTree()

	var m Int = 0
	var n Int = 1000
	for m < n {
		rbt.Insert(&RBNode{
			Parent: nil,
			Left:   nil,
			Right:  nil,
			color:  false,
			Item:   nil,
		})
		m++
	}
	if rbt.Len() != uint(n) {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), n)
	}

	for m > 0 {
		rbt.Delete(Int(m))
		m--
	}
	if rbt.Len() != 1 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 1)
	}
}
