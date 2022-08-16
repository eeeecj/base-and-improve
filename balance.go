package balance

import "fmt"

//  https://github.com/fibbery/go-balance/blob/master/balance/weight_round_robin_balance.go
type Balance interface {
	Balance(instances []*Instance, keys ...string) (*Instance, error)
}

type Manager struct {
	balances map[string]Balance
}

var manager = Manager{balances: make(map[string]Balance)}

func StartBalance(name string, instance []*Instance) (*Instance, error) {
	balance, ok := manager.balances[name]
	if !ok {
		err := fmt.Errorf("%s Balance Not Found", name)
		return nil, err
	}
	ins, err := balance.Balance(instance)
	return ins, err
}

func RegisterBalance(name string, b Balance) {
	manager.balances[name] = b
}
