package trie

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie(t *testing.T) {
	ti := NewTrie()
	ti.Add("xiaoming")
	ti.Add("xiaohua")
	ti.Add("datou")
	assert.Equal(t, true, ti.Search("xiaoming"))
	assert.Equal(t, false, ti.Remove("daxiong"))
	assert.Equal(t, true, ti.Remove("xiaoming"))
}
