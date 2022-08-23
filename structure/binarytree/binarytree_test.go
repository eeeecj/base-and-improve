package binarytree

import (
	"fmt"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	tree := NewBinaryTree()
	tree.Insert(1)
	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(4)
	tree.Insert(16)
	tree.Insert(5)
	tree.InorderTraversal(tree.Search(1), func(i int) {
		fmt.Println(i)
	})
	tree.PreorderTraversal(tree.Search(1), func(i int) {
		fmt.Println(i)
	})
	tree.PostorderTraversal(tree.Search(1), func(i int) {
		fmt.Println(i)
	})

}
