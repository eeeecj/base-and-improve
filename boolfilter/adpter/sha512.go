package adpter

import (
	"crypto/sha512"
	"hash"
)

type Sha512 struct {
	hash.Hash
}

func (m *Sha512) Sum64() uint64 {
	b := m.Sum(nil)
	return Base16ToUint64(b)
}

func NewSha512() hash.Hash64 {
	return &Sha512{sha512.New()}
}
