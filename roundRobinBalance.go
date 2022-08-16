package balance

import "fmt"

type RoundRobinBalance struct {
	index int
}

func init() {
	RegisterBalance("roundrobin", &RoundRobinBalance{})
}

func (r *RoundRobinBalance) Balance(instances []*Instance, keys ...string) (*Instance, error) {
	m := len(instances)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	if r.index >= m {
		r.index = 0
	}
	ins := instances[r.index]
	r.index = (r.index + 1) % m
	return ins, nil
}
