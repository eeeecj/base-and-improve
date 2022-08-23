package trie

type TrieImp interface {
	Add(s string)
	Search(s string) bool
	Remove(s string) bool
}

func NewTrie() TrieImp {
	return &Trie{root: &Node{children: make(map[int32]*Node)}}
}

type Node struct {
	isEnd    bool
	children map[int32]*Node
	Parent   *Node
}

type Trie struct {
	root *Node
}

func (t *Trie) Add(s string) {
	cur := t.root
	for _, v := range s {
		if _, ok := cur.children[v]; ok {
			cur = cur.children[v]
		} else {
			n := &Node{children: make(map[int32]*Node), Parent: cur}
			cur.children[v] = n
			cur = n
		}
	}
	cur.isEnd = true
}

func (t *Trie) Search(s string) bool {
	cur := t.root
	for _, v := range s {
		if _, ok := cur.children[v]; ok {
			cur = cur.children[v]
		} else {
			return false
		}
	}
	return cur.isEnd
}

func (t *Trie) Remove(s string) bool {
	cur := t.root
	for _, v := range s {
		if _, ok := cur.children[v]; ok {
			cur = cur.children[v]
		} else {
			return false
		}
	}
	if !cur.isEnd {
		return false
	}
	cur.isEnd = false
	index := len(s) - 1
	for len(cur.children) == 0 {
		delete(cur.children, int32(s[index]))
		cur = cur.Parent
		index--
	}
	return true
}
