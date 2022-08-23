package bitmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitmap(t *testing.T) {
	o := NewBitmap()
	o.Add(1264)
	o.Add(456)
	o.Add(984)
	o.Add(12)
	assert.Equal(t, true, o.Has(12))
	assert.Equal(t, false, o.Has(1))
	assert.Equal(t, 4, o.Len())
	o.Remove(456)
	assert.Equal(t, false, o.Has(456))
	assert.Equal(t, []int{12, 984, 1264}, o.Serialize())
}
