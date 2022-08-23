package linklist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoubleLinkList_Len(t *testing.T) {
	list := NewDoubleLinkList()
	for i := 0; i < 100; i++ {
		list.AddToTail(i)
	}
	assert.Equal(t, 100, list.Len())
}
func TestDoubleLinkList_Add(t *testing.T) {
	list := NewDoubleLinkList()
	list.Add(0, 0)
	list.Add(1, 2)
	list.Add(1, 1)
	assert.Equal(t, 2, list.Get(2).data)
}

func TestDoubleLinkList_AddToHeadTailGet(t *testing.T) {
	var nilnode *DoubleNode
	list := NewDoubleLinkList()
	list.AddToTail(0)
	list.AddToTail(1)
	list.Add(2, 2)
	list.Add(3, 3)
	list.Add(2, 5)
	assert.Equal(t, nilnode, list.Get(5))
	assert.Equal(t, 0, list.Get(0).data)
	assert.Equal(t, 3, list.Get(list.Len()-1).data)
	assert.Equal(t, 1, list.Get(1).data)
	assert.Equal(t, 2, list.Get(3).data)
	assert.Equal(t, 5, list.Get(2).data)
}

func TestDoubleLinkList_Delete(t *testing.T) {
	list := NewDoubleLinkList()
	assert.Equal(t, false, list.Delete(0))
	list.AddToHead(0)
	assert.Equal(t, true, list.Delete(0))
	list.AddToHead(0)
	list.AddToTail(1)
	list.AddToTail(2)
	list.AddToTail(3)
	list.AddToTail(4)
	assert.Equal(t, true, list.Delete(2))
	assert.Equal(t, true, list.Delete(0))
	assert.Equal(t, true, list.Delete(2))
	assert.Equal(t, []interface{}{1, 3}, list.ToList())
}
