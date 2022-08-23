package linklist

type LinkListImp interface {
	AddToFront(data interface{})
	AddToTail(data interface{})
	ToList() []interface{}
	Get(index int) *Node
	GetPrev(index int) *Node
	Add(index int, data interface{}) bool
	Delete(index int) bool
	Reverse() *LinkList
	IndexOf(data interface{}) int
	Len() int
}

type Node struct {
	data interface{}
	next *Node
}
type LinkList struct {
	head *Node
	size int
}

func NewLinkList() *LinkList {
	return &LinkList{head: &Node{}}
}

func (l *LinkList) Len() int {
	return l.size
}

func (l *LinkList) AddToFront(data interface{}) {
	l.size++
	l.head.next = &Node{data: data, next: l.head.next}
}

func (l *LinkList) AddToTail(data interface{}) {
	cur := l.head
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &Node{data: data}
	l.size++
}

func (l *LinkList) ToList() []interface{} {
	cur := l.head.next
	res := []interface{}{}
	for cur != nil {
		res = append(res, cur.data)
		cur = cur.next
	}
	return res
}

func (l *LinkList) Get(index int) *Node {
	if index > l.size || l.size == 0 {
		return nil
	}
	cur := l.head
	for index > 0 {
		cur = cur.next
		index--
	}
	return cur
}
func (l *LinkList) GetPrev(index int) *Node {
	if index <= 1 || l.size <= 1 || index >= l.size+1 {
		return nil
	}
	cur := l.head
	for index > 1 {
		cur = cur.next
		index--
	}
	return cur
}
func (l *LinkList) Add(index int, data interface{}) bool {
	if index > l.size+1 {
		return false
	}
	if index == 1 {
		l.AddToFront(data)
		return true
	}
	if index == l.size+1 {
		l.AddToTail(data)
		return true
	}
	prevNode := l.GetPrev(index)
	prevNode.next = &Node{data: data, next: prevNode.next}
	l.size++
	return true
}

func (l *LinkList) Delete(index int) bool {
	if index > l.size {
		return false
	}
	if index == 1 {
		l.head.next = l.head.next.next
		l.size--
		return true
	}
	prevNode := l.GetPrev(index)
	prevNode.next = prevNode.next.next
	l.size--
	return true
}

func (l *LinkList) Reverse() *LinkList {
	cur := l.head.next
	var pre *Node
	for cur != nil {
		next := cur.next
		cur.next = pre
		cur, pre = next, cur
	}
	l.head.next = pre
	return l
}

func (l *LinkList) IndexOf(data interface{}) int {
	if l.size == 0 {
		return -1
	}
	index := 1
	cur := l.head.next
	for cur != nil && cur.data != data {
		cur = cur.next
		index++
	}
	if cur == nil {
		return -1
	}
	return index
}
