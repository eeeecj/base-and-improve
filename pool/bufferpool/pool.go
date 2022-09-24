package bufferpool

import (
	"sort"
	"sync"
	"sync/atomic"
)

const (
	minBitSize = 6 // 2**6=64 is a CPU cache line size
	steps      = 20

	minSize = 1 << minBitSize
	maxSize = 1 << (minBitSize + steps - 1)

	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

type Pool struct {
	calls       [steps]uint64
	calibrating uint64
	defaultSize uint64
	maxSize     uint64
	pool        sync.Pool
}

var defaultPool Pool

func Get() *bufferPool { return defaultPool.Get() }
func (p *Pool) Get() *bufferPool {
	v := p.pool.Get()
	if v != nil {
		return v.(*bufferPool)
	}
	return &bufferPool{make([]byte, 0, p.defaultSize)}
}

func Put(b *bufferPool) { defaultPool.Put(b) }

func (p *Pool) Put(b *bufferPool) {
	index := indexOf(b.Len())
	if atomic.AddUint64(&p.calls[index], 1) > calibrateCallsThreshold {
		p.calibrate()
	}
	maxsize := int(atomic.LoadUint64(&p.maxSize))
	//如果maxSize等于0或容量小于maxSize的大小，则直接清空缓存
	//拒绝过大的内存造成内存浪费
	if maxsize == 0 || cap(b.B) <= maxsize {
		b.reset()
		p.pool.Put(b)
	}
}

func (p *Pool) calibrate() {
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}
	a := make(callSizes, 0, steps)
	var callSum uint64
	for i := uint64(0); i < steps; i++ {
		calls := atomic.SwapUint64(&p.calls[i], 0)
		callSum += calls
		a = append(a, callSize{
			calls: calls,
			size:  minSize << i,
		})
	}
	sort.Sort(a)
	defaultSize := a[0].size
	maxsize := defaultSize
	maxSum := uint64(float64(callSum) * maxPercentile)
	callSum = 0
	//从列表中获取调用次数最多的内存区间，使其可以满足95%的内存分配请求
	for i := 0; i < steps; i++ {
		if callSum > maxSum {
			break
		}
		callSum += a[i].calls
		size := a[i].size
		if size > maxsize {
			maxsize = size
		}
	}
	//将调用次数最多的内存大小设为默认值
	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxsize)
	atomic.StoreUint64(&p.calibrating, 0)
}

type callSize struct {
	calls uint64
	size  uint64
}
type callSizes []callSize

func (c callSizes) Less(i, j int) bool {
	return c[i].calls > c[j].calls
}

func (c callSizes) Len() int {
	return len(c)
}
func (c callSizes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func indexOf(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		idx = steps - 1
	}
	return steps
}
