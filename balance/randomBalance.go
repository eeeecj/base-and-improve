package balance

import (
	"fmt"
	"math/rand"
	"time"
)

type RandomBalance struct {
}

func init() {
	RegisterBalance("random", &RandomBalance{})
}

func (r *RandomBalance) Balance(ins []*Instance, keys ...string) (*Instance, error) {
	m := len(ins)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	rand.Seed(time.Now().UnixNano())
	return ins[rand.Intn(m)], nil
}
