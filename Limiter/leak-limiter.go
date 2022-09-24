package Limiter

import (
	"math"
	"sync"
	"time"
)

//漏桶算法是首先想象有一个木桶，桶的容量是固定的。
//当有请求到来时先放到木桶中，处理请求的worker以固定的速度从木桶中取出请求进行相应。
//如果木桶已经满了，直接返回请求频率超限的错误码或者页面。

//漏桶算法是流量最均匀的限流实现方式，一般用于流量“整形”。
//例如保护数据库的限流，先把对数据库的访问加入到木桶中，
//worker再以db能够承受的qps从木桶中取出请求，去访问数据库。

//木桶流入请求的速率是不固定的，但是流出的速率是恒定的。
//这样的话能保护系统资源不被打满，但是面对突发流量时会有大量请求失败，
//不适合电商抢购和微博出现热点事件等场景的限流。

type LeakLimiter struct {
	rate       float64
	capacity   float64
	water      float64
	lastLeakMs int64
	mu         sync.Mutex
}

func (ll *LeakLimiter) Allow() bool {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	leakWater := ll.water - (float64(now-ll.lastLeakMs) * ll.rate / 1000)
	ll.water = math.Max(0, leakWater)
	ll.lastLeakMs = now
	if ll.water+1 <= ll.capacity {
		ll.water++
		return true
	} else {
		return false
	}
}

func (ll *LeakLimiter) Set(rate float64, capacity float64) {
	ll.rate = rate
	ll.capacity = capacity
	ll.water = 0
	ll.lastLeakMs = time.Now().UnixNano() / 1e6
}
