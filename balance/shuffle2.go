package balance

import (
	"fmt"
	"math/rand"
	"time"
)

type Shuffle2 struct {
}

func init() {
	RegisterBalance("shuffle2", &Shuffle2{})
}

func (r *Shuffle2) Balance(ins []*Instance, keys ...string) (*Instance, error) {
	m := len(ins)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	rand.Seed(time.Now().UnixNano())
	for i := m - 1; i >= 0; i-- {
		end := i
		index := rand.Intn(i + 1)
		ins[index], ins[end] = ins[end], ins[index]
	}
	return ins[0], nil
}
