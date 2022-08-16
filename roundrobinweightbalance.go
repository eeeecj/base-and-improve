package balance

import "fmt"

type RoundRobinWightBalance struct {
	index  int
	weight int
}

func init() {
	RegisterBalance("roundrobinweight", NewRoundRobinWeightBalance())
}

func NewRoundRobinWeightBalance() *RoundRobinWightBalance {
	return &RoundRobinWightBalance{index: -1}
}

func (r *RoundRobinWightBalance) Balance(ins []*Instance, keys ...string) (*Instance, error) {
	m := len(ins)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	w := []int{}
	for i := 0; i < len(ins); i++ {
		w = append(w, ins[i].Weight)
	}
	g := getGcd(w)
	maxw := getMax(w)

	for {
		r.index = (r.index + 1) % m
		if r.index == 0 {
			r.weight = r.weight - g
			if r.weight <= 0 {
				r.weight = maxw
			}
		}
		instance := ins[r.index]
		if instance.Weight >= r.weight {
			return instance, nil
		}
	}
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
func getMax(nums []int) int {
	m := nums[0]
	for i := 1; i < len(nums); i++ {
		if m < nums[i] {
			m = nums[i]
		}
	}
	return m
}

func getGcd(nums []int) int {
	g := nums[0]
	for i := 0; i < len(nums); i++ {
		g = gcd(nums[i], g)
	}
	return g
}
