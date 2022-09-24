package Limiter

import (
	"sync"
	"time"
)

//计数器是一种比较简单粗暴的限流算法，其思想是在固定时间窗口内对请求进行计数，
//与阀值进行比较判断是否需要限流，一旦到了时间临界点，将计数器清零。
//计数器算法实现限流的问题是没有办法应对突发流量，不过它的算法实现起来确实最简单的，
//下面给出一个用Go代码实现的计数器。

type CountLimiter struct {
	rate  int
	begin time.Time
	cycle time.Duration
	count int
	lock  sync.Mutex
}

func (cl *CountLimiter) Allow() bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	if cl.count == cl.rate-1 {
		now := time.Now()
		if now.Sub(cl.begin) >= cl.cycle {
			cl.reset(now)
			return true
		}
		return false
	} else {
		cl.count++
		return true
	}
}
func (cl *CountLimiter) Set(rate int, cycle time.Duration) {
	cl.rate = rate
	cl.cycle = cycle
}
func (cl *CountLimiter) reset(now time.Time) {
	cl.begin = now
	cl.count = 0
}
