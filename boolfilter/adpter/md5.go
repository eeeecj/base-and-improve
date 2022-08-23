package adpter

import (
	"crypto/md5"
	"hash"
)

type MD5 struct {
	hash.Hash
}

func (m *MD5) Sum64() uint64 {
	b := m.Sum(nil)
	return Base16ToUint64(b)
}

func NewMD5() hash.Hash64 {
	return &MD5{md5.New()}
}
