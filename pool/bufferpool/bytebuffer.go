package bufferpool

import "io"

type bufferPool struct {
	B []byte
}

func (b *bufferPool) Len() int {
	return len(b.B)
}

func (b *bufferPool) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	start := int64(len(b.B))
	smax := int64(cap(b.B))
	n := start
	if smax == 0 {
		smax = 64
		p = make([]byte, smax)
	} else {
		p = b.B[:smax]
	}
	for {
		if n == smax {
			smax *= 2
			newb := make([]byte, smax)
			copy(newb, p)
			p = newb
		}
		nn, err := r.Read(p[start:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= start
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *bufferPool) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

func (b *bufferPool) Bytes() []byte {
	return b.B
}

func (b *bufferPool) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}
func (b *bufferPool) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

func (b *bufferPool) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

func (b *bufferPool) Set(p []byte) {
	b.B = append(b.B[:0], p...)
}
func (b *bufferPool) SetString(s string) {
	b.B = append(b.B[:0], s...)
}
func (b *bufferPool) reset() {
	b.B = b.B[:0]
}

func (b *bufferPool) String() string {
	return string(b.B)
}
