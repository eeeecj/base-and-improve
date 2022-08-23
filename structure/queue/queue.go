package queue

type QueueImp interface {
	Enqueue(data interface{})
	Dequeue() *Node
	Peek() interface{}
	IsEmpty() bool
	Serialize() []interface{}
	Len() int
}

type Queue struct {
	size int
	Head *Node
	Tail *Node
}

type Node struct {
	data interface{}
	next *Node
}

func NewQueue() QueueImp {
	return &Queue{}
}

func (q *Queue) Enqueue(data interface{}) {
	n := &Node{data: data}
	if q.Head == nil {
		q.Head = n
	} else {
		q.Tail.next = n
	}
	q.Tail = n
	q.size++
}

func (q *Queue) Dequeue() *Node {
	if q.size == 0 {
		return nil
	}
	head := q.Head
	q.Head = head.next
	q.size--
	head.next = nil
	return head
}

func (q *Queue) Peek() interface{} {
	if q.size == 0 {
		return nil
	}
	return q.Head.data
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}
func (q *Queue) Serialize() []interface{} {
	res := []interface{}{}
	cur := q.Head
	for cur != nil {
		res = append(res, cur.data)
		cur = cur.next
	}
	return res
}
func (q *Queue) Len() int {
	return q.size
}
