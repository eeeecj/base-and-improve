package grpc_pool

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

type Pool interface {
	Get() (ConnImp, error)
	Close() error
	Status() string
}
type pool struct {
	index   uint32
	current int32
	ref     int32
	opt     Options
	conns   []*conn
	address string
	closed  int32
	sync.RWMutex
}

func New(addr string, opt Options) (Pool, error) {
	if addr == "" {
		return nil, errors.New("invalid address settings")
	}
	if opt.Dial == nil {
		return nil, errors.New("invalid dial settings")
	}
	if opt.MaxIdle <= 0 || opt.MaxActive <= 0 || opt.MaxIdle > opt.MaxActive {
		return nil, errors.New("invalid maximum settings")
	}
	if opt.MaxConcurrentStreams <= 0 {
		return nil, errors.New("invalid maximun settings")
	}
	p := &pool{
		index:   0,
		current: int32(opt.MaxIdle),
		ref:     0,
		opt:     opt,
		conns:   make([]*conn, opt.MaxActive),
		address: addr,
		closed:  0,
	}

	for i := 0; i < p.opt.MaxIdle; i++ {
		c, err := p.opt.Dial(addr)
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("dial is not able to fill the pool: %s", err)
		}
		p.conns[i] = p.wrapConn(c, false)
	}
	return p, nil
}

func (p *pool) incrRef() int32 {
	newRef := atomic.AddInt32(&p.ref, 1)
	if newRef == math.MaxInt32 {
		panic(fmt.Sprintf("overflow ref: %d", newRef))
	}
	return newRef
}

func (p *pool) decrRef() {
	newRef := atomic.AddInt32(&p.ref, -1)
	if newRef < 0 && atomic.LoadInt32(&p.closed) == 0 {
		panic(fmt.Sprintf("negative ref: %d", newRef))
	}
	if newRef == 0 && atomic.LoadInt32(&p.current) > int32(p.opt.MaxIdle) {
		p.Lock()
		if atomic.LoadInt32(&p.ref) == 0 {
			atomic.StoreInt32(&p.current, int32(p.opt.MaxIdle))
			p.deleteFrom(p.opt.MaxIdle)
		}
		p.Unlock()
	}
}
func (p *pool) reset(index int) {
	con := p.conns[index]
	if con == nil {
		return
	}
	con.reset()
	p.conns[index] = nil
}
func (p *pool) deleteFrom(start int) {
	for i := start; i < p.opt.MaxActive; i++ {
		p.reset(i)
	}
}

func (p *pool) Get() (ConnImp, error) {
	nextR := p.incrRef()
	p.RLock()
	cur := atomic.LoadInt32(&p.current)
	p.RUnlock()
	if p.current == 0 {
		return nil, errors.New("pool is closed")
	}
	if nextR <= cur*int32(p.opt.MaxConcurrentStreams) {
		next := atomic.AddUint32(&p.index, 1) % uint32(cur)
		return p.conns[next], nil
	}
	if cur == int32(p.opt.MaxActive) {
		if p.opt.Reuse {
			next := atomic.AddUint32(&p.index, 1) % uint32(cur)
			return p.conns[next], nil
		}
		c, err := p.opt.Dial(p.address)
		return p.wrapConn(c, false), err
	}
	p.Lock()
	cur = atomic.LoadInt32(&p.current)
	if cur < int32(p.opt.MaxActive) && nextR > cur*int32(p.opt.MaxConcurrentStreams) {
		incr := cur
		if cur+incr > int32(p.opt.MaxActive) {
			incr = int32(p.opt.MaxActive) - cur
		}
		var i int32
		var err error
		for i = 0; i < incr; i++ {
			c, er := p.opt.Dial(p.address)
			if er != nil {
				err = er
				break
			}
			p.reset(int(cur + i))
			p.conns[cur+i] = p.wrapConn(c, false)
		}
		cur += i
		atomic.StoreInt32(&p.current, cur)
		if err != nil {
			p.Unlock()
			return nil, err
		}
	}
	p.Unlock()
	next := atomic.AddUint32(&p.index, 1) % uint32(cur)
	return p.conns[next], nil
}
func (p *pool) Close() error {
	atomic.StoreInt32(&p.ref, 0)
	atomic.StoreInt32(&p.current, 0)
	atomic.StoreUint32(&p.index, 0)
	atomic.StoreInt32(&p.closed, 1)
	p.deleteFrom(0)
	return nil
}

func (p *pool) Status() string {
	return fmt.Sprintf("address:%s, index:%d, current:%d, ref:%d. option:%v",
		p.address, p.index, p.current, p.ref, p.opt)
}
