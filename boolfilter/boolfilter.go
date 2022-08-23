package boolfilter

import (
	"github.com/balance/boolfilter/adpter"
	"github.com/balance/boolfilter/global"
	"hash"
)

var DefaultHash = []global.HashFunc{
	func() hash.Hash64 { return adpter.NewMD5() },
	func() hash.Hash64 { return adpter.NewSha1() },
	func() hash.Hash64 { return adpter.NewSha256() },
	func() hash.Hash64 { return adpter.NewSha512() },
}

type Adapter interface {
	Clear() error
	Push(content []byte) error
	Write() error
	Exists(content []byte) (bool, error)
	IsEmpty() bool
	Close() error
}
