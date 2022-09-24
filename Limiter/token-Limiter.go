package Limiter

import (
	"sync"
	"time"
)

//令牌桶是反向的"漏桶"，它是以恒定的速度往木桶里加入令牌，
//木桶满了则不再加入令牌。服务收到请求时尝试从木桶中取出一个令牌，
//如果能够得到令牌则继续执行后续的业务逻辑。
//如果没有得到令牌，直接返回访问频率超限的错误码或页面等，不继续执行后续的业务逻辑。

//特点：由于木桶内只要有令牌，请求就可以被处理，所以令牌桶算法可以支持突发流量。

//同时由于往木桶添加令牌的速度是恒定的，且木桶的容量有上限，
//所以单位时间内处理的请求书也能够得到控制，起到限流的目的

//适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。

type TokenLimiter struct {
	rate         int64
	capacity     int64
	tokens       int64
	lastTokenSec int64
	mu           sync.Mutex
}

func (tl *TokenLimiter) Allow() bool {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	now := time.Now().Unix()
	tl.tokens = tl.tokens + (now-tl.lastTokenSec)*tl.rate
	if tl.tokens > tl.capacity {
		tl.tokens = tl.capacity
	}

	tl.lastTokenSec = now
	if tl.tokens > 0 {
		tl.tokens--
		return true
	}
	return false
}

func (tl *TokenLimiter) Set(rate int64, capacity int64) {
	tl.rate = rate
	tl.capacity = capacity
	tl.tokens = 0
	tl.lastTokenSec = time.Now().Unix()
}
