package hashtable

import (
	"hash/fnv"
)

const MaxTableSize = 256

type Item struct {
	key, value interface{}
	next       *Item
}

type HashTableImp interface {
	Add(key, value interface{})
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Remove(key interface{})
}

type HashTable struct {
	data [MaxTableSize]*Item
}

func NewHashTable() HashTableImp {
	return &HashTable{}
}

func (h *HashTable) Add(key, value interface{}) {
	pos := generateHash(key)
	if h.data[pos] == nil {
		h.data[pos] = &Item{key: key, value: value}
		return
	}
	cur := h.data[pos]
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &Item{key: key, value: value}
}
func (h *HashTable) Get(key interface{}) interface{} {
	pos := generateHash(key)
	cur := h.data[pos]
	for cur != nil {
		if cur.key == key {
			return cur.value
		}
		cur = cur.next
	}
	return nil
}
func (h *HashTable) Set(key, value interface{}) {
	pos := generateHash(key)
	cur := h.data[pos]
	if h.data[pos] == nil {
		h.data[pos] = &Item{key: key, value: value}
		return
	}
	for cur.next != nil {
		if cur.key == key {
			cur.value = value
			return
		}
		cur = cur.next
	}
	if cur.key == key {
		cur.value = value
		return
	}
	cur.next = &Item{key: key, value: value}
}

func (h *HashTable) Remove(key interface{}) {
	pos := generateHash(key)
	pre := &Item{next: h.data[pos]}
	n := pre
	cur := pre.next
	for cur != nil {
		if cur.key == key {
			pre.next = cur.next
			h.data[pos] = n
			return
		}
		cur, pre = cur.next, pre.next
	}
}
func generateHash(key interface{}) uint8 {
	hash := fnv.New32a()
	hash.Write([]byte(key.(string)))
	return uint8(hash.Sum32() % MaxTableSize)
}
