package linklist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkList_Len(t *testing.T) {
	list := LinkListImp(NewLinkList())
	list.AddToFront(1)
	assert.Equal(t, 1, list.Len())
	for i := 0; i < 99; i++ {
		list.AddToTail(1)
	}
	assert.Equal(t, 100, list.Len())
}

func TestLinkList_AddToFrontAndTail(t *testing.T) {
	list := LinkListImp(NewLinkList())
	list.AddToFront(1)
	assert.Equal(t, 1, list.Get(1).data)
	list.AddToTail(2)
	assert.Equal(t, 2, list.Get(2).data)
	assert.Equal(t, []interface{}{1, 2}, list.ToList())
}
func TestLinkList_AddAndDelete(t *testing.T) {
	var nilnode *Node
	list := LinkListImp(NewLinkList())
	assert.Equal(t, -1, list.IndexOf(1))
	list.AddToFront(1)
	list.AddToTail(2)
	list.Add(1, 3)
	assert.Equal(t, 3, list.Get(1).data)
	assert.Equal(t, 1, list.Get(2).data)
	list.Delete(2)
	assert.Equal(t, []interface{}{3, 2}, list.ToList())
	assert.Equal(t, 2, list.Get(2).data)
	assert.Equal(t, 2, list.Len())
	assert.Equal(t, nilnode, list.Get(4))
	assert.Equal(t, nilnode, list.GetPrev(1))
	assert.Equal(t, false, list.Add(4, 4))
	assert.Equal(t, true, list.Add(3, 3))
	assert.Equal(t, true, list.Add(3, 4))
	assert.Equal(t, 4, list.Get(3).data)
	assert.Equal(t, false, list.Delete(6))
	assert.Equal(t, true, list.Delete(1))
	assert.Equal(t, -1, list.IndexOf(6))
	assert.Equal(t, 1, list.IndexOf(2))
}

func TestLinkList_Reverse(t *testing.T) {
	list := LinkListImp(NewLinkList())
	list.AddToFront(1)
	list.AddToTail(2)
	list.Add(1, 3)
	list.Reverse()
	assert.Equal(t, 2, list.Get(1).data)
	assert.Equal(t, 3, list.Get(3).data)
}
