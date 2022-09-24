package Limiter

import (
	"sync"
	"time"
)

//滑动窗口算法将一个大的时间窗口分成多个小窗口，每次大窗口向后滑动一个小窗口，
//并保证大的窗口内流量不会超出最大值，这种实现比固定窗口的流量曲线更加平滑。

//对于滑动时间窗口，我们可以把1ms的时间窗口划分成10个小窗口，
//或者想象窗口有10个时间插槽time slot, 每个time slot统计某个100ms的请求数量。
//每经过100ms，有一个新的time slot加入窗口，
//早于当前时间1s的time slot出窗口。窗口内最多维护10个time slot。

//滑动窗口算法是固定窗口的一种改进，
//但从根本上并没有真正解决固定窗口算法的临界突发流量问题

type timeSlot struct {
	timestamp time.Time
	count     int
}

func countReq(win []*timeSlot) int {
	var count int
	for _, ts := range win {
		count += ts.count
	}
	return count
}

type SlideWinLimiter struct {
	mu           sync.Mutex
	SlotDuration time.Duration
	WinDuration  time.Duration
	numSlots     int
	wins         []*timeSlot
	maxReq       int
}

func NewSlideWinLimiter(sd time.Duration, wd time.Duration, maxreq int) *SlideWinLimiter {
	return &SlideWinLimiter{
		SlotDuration: sd,
		WinDuration:  wd,
		numSlots:     int(wd / sd),
		maxReq:       maxreq,
	}
}
func (sl *SlideWinLimiter) Allow() bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	now := time.Now()
	timeoutOffset := -1
	for i, ts := range sl.wins {
		if ts.timestamp.Add(sl.WinDuration).After(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		sl.wins = sl.wins[timeoutOffset+1:]
	}
	var res bool
	if sl.maxReq > countReq(sl.wins) {
		res = true
	}
	if len(sl.wins) > 0 {
		last := sl.wins[len(sl.wins)-1]
		if last.timestamp.Add(sl.SlotDuration).Before(now) {
			last = &timeSlot{timestamp: now, count: 1}
			sl.wins = append(sl.wins, last)
		} else {
			last.count++
		}
	} else {
		last := &timeSlot{timestamp: now, count: 1}
		sl.wins = append(sl.wins, last)
	}
	return res
}
