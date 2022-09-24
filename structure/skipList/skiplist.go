package skipList

import (
	"math/bits"
	"math/rand"
	"sync"
	"time"
)

var defaultMaxLevel = 10

type Node struct {
	next []*Element
}
type Element struct {
	Node
	key   interface{}
	value interface{}
}

type SkipList struct {
	mu           *sync.Mutex
	head         Node
	maxLevel     int
	len          int
	preNodeCache []*Node
	rander       *rand.Rand
	keyCmp       func(a, b interface{}) int
}

func New() *SkipList {
	l := &SkipList{
		mu:       &sync.Mutex{},
		maxLevel: defaultMaxLevel,
		rander:   rand.New(rand.NewSource(time.Now().Unix())),
	}
	l.head.next = make([]*Element, l.maxLevel)
	l.preNodeCache = make([]*Node, l.maxLevel)
	return l
}
func (s *SkipList) Insert(key, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prevs := s.findPrevNodes(key)
	//如果插入的值已经存在
	if prevs[0].next[0] != nil && s.keyCmp(prevs[0].next[0].key, key) == 0 {
		prevs[0].next[0].value = value
		return
	}
	level := s.randomLevel()
	e := &Element{
		Node:  Node{next: make([]*Element, level)},
		key:   key,
		value: value,
	}
	for i := range e.next {
		//将新插入的点指向prev的next，并将prev的next指向新插入节点
		e.next[i] = prevs[i].next[i]
		prevs[i].next[i] = e
	}
	s.len++
}

func (s *SkipList) Get(key interface{}) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	var pre = &s.head
	for i := s.maxLevel - 1; i >= 0; i-- {
		cur := pre.next[i]
		for ; cur != nil; cur = cur.next[i] {
			cmpRet := s.keyCmp(cur.key, key)
			if cmpRet == 0 {
				return cur.value
			}
			if cmpRet > 0 {
				break
			}
			pre = &cur.Node
		}
	}
	return nil
}

func (s *SkipList) Remove(key interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	prevs := s.findPrevNodes(key)
	element := prevs[0].next[0]
	if element == nil {
		return false
	}
	if element != nil && s.keyCmp(element.key, key) != 0 {
		return false
	}
	for i, v := range element.next {
		prevs[i].next[i] = v
	}
	s.len--
	return true
}
func (s *SkipList) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.len
}

func (s *SkipList) randomLevel() int {
	total := uint64(1)<<uint64(s.maxLevel) - 1
	k := s.rander.Uint64() % total
	level := s.maxLevel - bits.Len64(k) + 1
	for level > 1 && 1<<(level-3) > s.len {
		level--
	}
	return level
}
func (s *SkipList) findPrevNodes(key interface{}) []*Node {
	prevs := s.preNodeCache
	prev := &s.head
	for i := s.maxLevel - 1; i >= 0; i-- {
		if s.head.next[i] != nil {
			for next := prev.next[i]; next != nil; next = next.next[i] {
				if s.keyCmp(next.key, key) >= 0 {
					break
				}
				prev = &next.Node
			}
		}
		prevs[i] = prev
	}
	return prevs
}
func (s *SkipList) Traversal(visitor func(k, v interface{}) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for e := s.head.next[0]; e != nil; e = e.next[0] {
		if !visitor(e.key, e.value) {
			return
		}
	}
}

func (s *SkipList) Keys() []interface{} {
	var keys []interface{}
	s.Traversal(func(k, v interface{}) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}
