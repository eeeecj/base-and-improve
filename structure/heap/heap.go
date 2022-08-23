package heap

type HeapImp interface {
	Add(data int)
	AddSlice(data []int)
	Len() int
	BuildMaxHeap(n int) []int
	BuildMinHeap(n int) []int
	Sort() []int
	SortDesc() []int
}

func NewHeap() HeapImp {
	return &Heap{}
}

type Heap struct {
	data []int
}

func (h *Heap) Add(data int) {
	h.data = append(h.data, data)
}

func (h *Heap) AddSlice(data []int) {
	h.data = append(h.data, data...)
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) heapfy_max(i int, n int) {
	left := i*2 + 1
	right := i*2 + 2
	large := i
	if left <= n && h.data[left] > h.data[large] {
		large = left
	}
	if right <= n && h.data[right] > h.data[large] {
		large = right
	}
	if large != i {
		h.data[i], h.data[large] = h.data[large], h.data[i]
		h.heapfy_max(large, n)
	}
}

func (h *Heap) BuildMaxHeap(n int) []int {
	for i := n / 2; i >= 0; i-- {
		h.heapfy_max(i, n)
	}
	return h.data
}

func (h *Heap) BuildMinHeap(n int) []int {
	for i := n / 2; i >= 0; i-- {
		h.heapfy_min(i, n)
	}
	return h.data
}

func (h *Heap) heapfy_min(i int, n int) {
	left := i*2 + 1
	right := i*2 + 2
	least := i
	if left <= n && h.data[left] < h.data[least] {
		least = left
	}
	if right <= n && h.data[right] < h.data[least] {
		least = right
	}
	if least != i {
		h.data[i], h.data[least] = h.data[least], h.data[i]
		h.heapfy_min(least, n)
	}
}

func (h *Heap) Sort() []int {
	l := h.Len() - 1
	for l >= 0 {
		h.BuildMaxHeap(l)
		h.data[0], h.data[l] = h.data[l], h.data[0]
		l--
	}
	return h.data
}

func (h *Heap) SortDesc() []int {
	l := h.Len() - 1
	for l >= 1 {
		h.BuildMinHeap(l)
		h.data[0], h.data[l] = h.data[l], h.data[0]
		l--
	}
	return h.data
}
