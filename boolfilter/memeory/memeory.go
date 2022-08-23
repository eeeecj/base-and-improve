package memeory

import (
	"errors"
	"github.com/balance/boolfilter"
	"github.com/balance/boolfilter/global"
)

func NewMemoryFilter(content []byte, hashes ...global.HashFunc) boolfilter.Adapter {
	return &Filter{
		Bits:      content,
		Hashes:    hashes,
		IsChanged: false,
	}
}

type Filter struct {
	Bits      []byte
	Hashes    []global.HashFunc
	IsChanged bool
}

func (f *Filter) Clear() error {
	for i, _ := range f.Bits {
		f.Bits[i] = 0
	}
	return nil
}
func (f *Filter) Push(content []byte) error {
	var byteLen = uint64(len(f.Bits))
	if byteLen < 1 {
		return errors.New("Empty Content")
	}
	for _, h := range f.Hashes {
		v := h()
		v.Reset()
		_, err := v.Write(content)
		if err != nil {
			return err
		}
		res := v.Sum64()
		word, bit := res%byteLen, (res/byteLen)&7
		now := f.Bits[word] | 1<<bit
		if now != f.Bits[word] {
			f.Bits[word] = now
			f.IsChanged = true
		}
	}
	return nil
}

func (f *Filter) Write() error {
	f.IsChanged = false
	return nil
}

func (f *Filter) Exists(content []byte) (bool, error) {
	blen := uint64(len(f.Bits))
	if blen < 1 {
		return false, errors.New("Empty Content")
	}
	for _, h := range f.Hashes {
		v := h()
		v.Reset()
		_, err := v.Write(content)
		if err != nil {
			return false, err
		}
		res := v.Sum64()
		word, bit := res%blen, (res/blen)&7
		if f.Bits[word]|1<<bit != f.Bits[word] {
			return false, nil
		}
	}
	return true, nil
}
func (f *Filter) IsEmpty() bool {
	for i, _ := range f.Bits {
		if f.Bits[i] != 0 {
			return false
		}
	}
	return true
}

func (f *Filter) Close() error {
	f.IsChanged = false
	return nil
}
