package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// If workes idle for at least this period of time, then stop a worker.
	idleTimeout = 2 * time.Second
)

type WorkerPool struct {
	maxWorkers   int
	taskQueue    chan func()
	workerQueue  chan func()
	stoppedChan  chan struct{}
	stopSignal   chan struct{}
	waitingQueue []func()
	stopLock     sync.Mutex
	stopOnce     sync.Once
	stopped      bool
	waiting      int32
	wait         bool
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	pool := &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func()),
		workerQueue: make(chan func()),
		stoppedChan: make(chan struct{}),
		stopSignal:  make(chan struct{}),
	}
	go pool.dispatch()
	return pool
}

func (p *WorkerPool) Size() int {
	return p.maxWorkers
}
func (p *WorkerPool) Stop() {
	p.stop(false)
}

func (p *WorkerPool) StopWait() {
	p.stop(true)
}
func (p *WorkerPool) Stopped() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	return p.stopped
}
func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
	return
}

func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func() {
		task()
		close(doneChan)
	}
	<-doneChan
}
func (p *WorkerPool) WaitingSize() int {
	return len(p.waitingQueue)
}

func (p *WorkerPool) Pause(ctx context.Context) {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	if p.stopped {
		return
	}
	ready := new(sync.WaitGroup)
	ready.Add(p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		p.Submit(func() {
			ready.Done()
			select {
			case <-ctx.Done():
			case <-p.stopSignal:
			}
		})
	}
	ready.Wait()
}
func (p *WorkerPool) stop(wait bool) {
	p.stopOnce.Do(func() {
		close(p.stopSignal)
		p.stopLock.Lock()
		p.stopped = true
		p.stopLock.Unlock()
		p.wait = wait
		close(p.taskQueue)
	})
	<-p.stoppedChan
}
func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	time := time.NewTimer(idleTimeout)
	var workerCount int
	var idle bool
	var wg sync.WaitGroup

Loop:
	for {
		if len(p.waitingQueue) != 0 {
			if !p.processWaitingPool() {
				break Loop
			}
			continue
		}
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				break Loop
			}
			select {
			case p.workerQueue <- task:
			default:
				if workerCount < p.maxWorkers {
					wg.Add(1)
					go worker(task, p.workerQueue, &wg)
					workerCount++
				} else {
					p.waitingQueue = append(p.waitingQueue, task)
					atomic.StoreInt32(&p.waiting, int32(len(p.waitingQueue)))
				}
			}
			idle = false
		case <-time.C:
			if idle && workerCount > 0 {
				if p.killIdleWorker() {
					workerCount--
				}
			}
			idle = true
			time.Reset(idleTimeout)
		}
	}
	if p.wait {
		p.runQueueTask()
	}
	for workerCount > 0 {
		p.workerQueue <- nil
		workerCount--
	}
	wg.Wait()
	time.Stop()
}
func (p *WorkerPool) runQueueTask() {
	for len(p.waitingQueue) != 0 {
		p.workerQueue <- p.waitingQueue[0]
		p.waitingQueue = p.waitingQueue[1:]
		atomic.StoreInt32(&p.waiting, int32(len(p.waitingQueue)))
	}
}
func (p *WorkerPool) killIdleWorker() bool {
	select {
	case p.workerQueue <- nil:
		return true
	default:
		return false
	}
}
func worker(task func(), workerQueue chan func(), wg *sync.WaitGroup) {
	for task != nil {
		task()
		task = <-workerQueue
	}
	wg.Done()
}
func (p *WorkerPool) processWaitingPool() bool {
	select {
	case task, ok := <-p.taskQueue:
		if !ok {
			return false
		}
		p.waitingQueue = append(p.waitingQueue, task)
	case p.workerQueue <- p.waitingQueue[0]:
		p.waitingQueue = p.waitingQueue[1:]
	}
	atomic.StoreInt32(&p.waiting, int32(len(p.waitingQueue)))
	return true
}
