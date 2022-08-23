package stack

type StackImp interface {
	Push(data interface{})
	Pop() *Node
	Peek() interface{}
	IsEmpty() bool
	Serialize() []interface{}
	Len() int
}
type Stack struct {
	Tail *Node
	size int
}
type Node struct {
	data interface{}
	next *Node
}

func NewStack() StackImp {
	return &Stack{}
}
func (s *Stack) Push(data interface{}) {
	n := &Node{data: data}
	if s.Tail != nil {
		n.next = s.Tail
	}
	s.Tail = n
	s.size++
}
func (s *Stack) Pop() *Node {
	if s.size == 0 {
		return nil
	}
	tail := s.Tail
	s.Tail = tail.next
	tail.next = nil
	s.size--
	return tail
}

func (s *Stack) Peek() interface{} {
	if s.size == 0 {
		return nil
	}
	return s.Tail.data
}

func (s *Stack) IsEmpty() bool {
	return s.size == 0
}
func (s *Stack) Serialize() []interface{} {
	cur := s.Tail
	res := []interface{}{}
	for cur != nil {
		res = append(res, cur.data)
		cur = cur.next
	}
	return res
}
func (s *Stack) Len() int {
	return s.size
}
