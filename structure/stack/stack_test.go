package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(4)
	s.Push(8)
	assert.Equal(t, []interface{}{8, 4, 2, 1}, s.Serialize())
	assert.Equal(t, 8, s.Pop().data)
	assert.Equal(t, 4, s.Peek())
	s.Pop()
	assert.Equal(t, 2, s.Len())
	assert.Equal(t, []interface{}{2, 1}, s.Serialize())
}
