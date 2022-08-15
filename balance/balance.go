package balance

type Balance interface {
	Balance([]*Instance, keys ...string) (*Instance, error)
}
