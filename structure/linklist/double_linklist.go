package linklist

type DoubleLinkList struct {
	Head *DoubleNode
	Tail *DoubleNode
	size int
}

type DoubleNode struct {
	data interface{}
	prev *DoubleNode
	next *DoubleNode
}

func NewDoubleLinkList() *DoubleLinkList {
	return &DoubleLinkList{}
}

func (d *DoubleLinkList) Len() int {
	return d.size
}

func (d *DoubleLinkList) Get(index int) *DoubleNode {
	if d.size == 0 || d.size < index+1 || index < 0 {
		return nil
	}
	if index == 0 {
		return d.Head
	} else if index == d.size-1 {
		return d.Tail
	}
	if index*2 > d.size-1 {
		cur := d.Tail
		idx := d.size - index - 1
		for idx > 0 {
			cur = cur.prev
			idx--
		}
		return cur
	} else {
		cur := d.Head
		idx := index
		for idx > 0 {
			cur = cur.next
			idx--
		}
		return cur
	}
}

func (d *DoubleLinkList) Add(index int, data interface{}) {
	if index == 0 {
		d.AddToHead(data)
		return
	}
	if index == d.size {
		d.AddToTail(data)
		return
	}
	idxnode := d.Get(index)
	prev := idxnode.prev
	node := &DoubleNode{data: data, prev: prev, next: idxnode}
	prev.next = node
	idxnode.prev = node
	d.size++
}

func (d *DoubleLinkList) AddToHead(data interface{}) {
	node := &DoubleNode{data: data}
	if d.size == 0 {
		d.Head = node
		d.Tail = node
	} else {
		node.next = d.Head
		d.Head.prev = node
		d.Head = node
	}
	d.size++
}
func (d *DoubleLinkList) AddToTail(data interface{}) {
	node := &DoubleNode{data: data}
	if d.size == 0 {
		d.Head = node
		d.Tail = node
	} else {
		node.prev = d.Tail
		d.Tail.next = node
		d.Tail = node
	}
	d.size++
}

func (d *DoubleLinkList) Delete(index int) bool {
	idxnode := d.Get(index)
	if idxnode == nil {
		return false
	}
	if idxnode == d.Head && idxnode == d.Tail {
		d.Head = nil
		d.Tail = nil
		d.size--
		return true
	}
	if idxnode == d.Head {
		d.Head = d.Head.next
		d.Head.prev = nil
		d.size--
		return true
	}
	if idxnode == d.Tail {
		d.Tail = d.Tail.prev
		d.Tail.next = nil
		d.size--
		return true
	}
	prev := idxnode.prev
	next := idxnode.next
	prev.next = next
	next.prev = prev
	idxnode = nil
	d.size--
	return true
}

func (d *DoubleLinkList) ToList() []interface{} {
	cur := d.Head
	res := []interface{}{}
	for cur != nil {
		res = append(res, cur.data)
		cur = cur.next
	}
	return res
}
