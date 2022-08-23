package bitmap

type BitmapImp interface {
	Has(num int) bool
	Add(num int)
	Remove(num int)
	Len() int
	Serialize() []int
}

func NewBitmap() BitmapImp {
	return &Bitmap{}
}

type Bitmap struct {
	bits []uint64
	size int
}

func (b *Bitmap) Has(num int) bool {
	word, bit := num/64, uint(num%64)
	return len(b.bits) > word && (b.bits[word])&(1<<bit) != 0
}

func (b *Bitmap) Add(num int) {
	word, bit := num/64, uint(num%64)
	td := word - len(b.bits) + 1
	if td > 0 {
		b.bits = append(b.bits, make([]uint64, td)...)
	}
	if !b.Has(num) {
		b.bits[word] |= 1 << bit
		b.size++
	}
}
func (b *Bitmap) Remove(num int) {
	word, bit := num/64, uint(num%64)
	if b.Has(num) {
		b.bits[word] ^= 1 << bit
		b.size--
	}
}

func (b *Bitmap) Len() int {
	return b.size
}
func (b *Bitmap) Serialize() []int {
	res := []int{}
	base := 64
	for i := 0; i < len(b.bits); i++ {
		t := b.bits[i]
		index := 0
		for t > 0 {
			if t&1 == 1 {
				res = append(res, index+base*i)
			}
			t >>= 1
			index++
		}
	}
	return res
}
