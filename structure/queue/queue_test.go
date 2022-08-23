package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue(t *testing.T) {
	var q = NewQueue()
	q.Enqueue(1)
	q.Enqueue(4)
	q.Enqueue(8)
	q.Enqueue(2)

	assert.Equal(t, []interface{}{1, 4, 8, 2}, q.Serialize())
	//return
	assert.Equal(t, 1, q.Peek())

	q.Dequeue()
	assert.Equal(t, 4, q.Dequeue().data)

	assert.Equal(t, 2, q.Len())

	assert.Equal(t, []interface{}{8, 2}, q.Serialize())
	assert.Equal(t, false, q.IsEmpty())
}
