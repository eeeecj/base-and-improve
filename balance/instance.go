package balance

type Instance struct {
	Host   string
	Port   int
	Weight int
}

func NewInstance(host string, port int, w int) *Instance {
	return &Instance{
		Host:   host,
		Port:   port,
		Weight: w,
	}
}
