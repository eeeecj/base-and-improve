package balance

import (
	"fmt"
	"math/rand"
	"time"
)

type Shuffle struct {
}

func init() {
	RegisterBalance("shuffle", &Shuffle{})
}
func (r *Shuffle) Balance(ins []*Instance, keys ...string) (*Instance, error) {
	m := len(ins)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i <= m/2; i++ {
		a := rand.Intn(m)
		b := rand.Intn(m)
		ins[a], ins[b] = ins[b], ins[a]
	}
	return ins[rand.Perm(m)[0]], nil
}
