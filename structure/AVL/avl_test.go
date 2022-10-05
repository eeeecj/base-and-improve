package AVL

import (
	"fmt"
	"testing"
)

func TestAVL(t *testing.T) {
	avlTree := NewAVLTree(3)
	// 插入节点
	avlTree.Insert(2)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(1)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(4)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(5)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(6)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(7)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(10)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(9)
	avlTree.Traverse()
	fmt.Println()
	avlTree.Insert(8)
	avlTree.Traverse()
	fmt.Println()
	// 判断是否是平衡二叉树
	fmt.Print("avlTree 是平衡二叉树: ")
	fmt.Println(avlTree.IsAVLTree())
	// 中序遍历生成的二叉树看是否是二叉排序树
	fmt.Print("中序遍历结果: ")
	avlTree.Traverse()
	fmt.Println()
	// 查找节点
	fmt.Print("查找值为 5 的节点: ")
	fmt.Printf("%v\n", avlTree.Find(5))
	// 删除节点
	avlTree.Delete(5)
	// 删除后是否还是平衡二叉树
	fmt.Print("avlTree 仍然是平衡二叉树: ")
	fmt.Println(avlTree.IsAVLTree())
	fmt.Print("删除节点后的中序遍历结果: ")
	avlTree.Traverse()
	fmt.Println()
}
