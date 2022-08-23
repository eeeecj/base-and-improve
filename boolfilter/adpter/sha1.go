package adpter

import (
	"crypto/sha1"
	"hash"
)

type Sha1 struct {
	hash.Hash
}

func (m *Sha1) Sum64() uint64 {
	b := m.Sum(nil)
	return Base16ToUint64(b)
}

func NewSha1() hash.Hash64 {
	return &Sha1{sha1.New()}
}
