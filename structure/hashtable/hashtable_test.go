package hashtable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashTable(t *testing.T) {
	h := NewHashTable()
	h.Add("xiaoming", 12)
	h.Add("xiaohua", 15)
	h.Add("adong", 1)
	h.Add("abo", 14)
	assert.Equal(t, 12, h.Get("xiaoming").(int))
	assert.Equal(t, 1, h.Get("adong").(int))
	h.Set("xiaoming", 129)
	assert.Equal(t, 129, h.Get("xiaoming").(int))
	h.Remove("xiaoming")
	assert.Equal(t, nil, h.Get("xiaoming"))
}
