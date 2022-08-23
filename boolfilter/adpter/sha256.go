package adpter

import (
	"crypto/sha256"
	"hash"
)

type Sha256 struct {
	hash.Hash
}

func (m *Sha256) Sum64() uint64 {
	b := m.Sum(nil)
	return Base16ToUint64(b)
}

func NewSha256() hash.Hash64 {
	return &Sha256{sha256.New()}
}
