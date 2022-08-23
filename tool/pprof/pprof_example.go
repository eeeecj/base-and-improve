//package main
//
//import (
//	"net/http"
//	_ "net/http/pprof"
//	"sync"
//)
//
//var m = &sync.Mutex{}
//
//func main() {
//	go func() {
//		m.Lock()
//		go func() {
//			m.Lock()
//			m.Unlock()
//		}()
//		m.Unlock()
//	}()
//	http.ListenAndServe(":6060", nil)
//}
//
//var datas []string
//
//func Add(str string) string {
//	data := []byte(str)
//	sData := string(data)
//	datas = append(datas, sData)
//
//	return sData
//}

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

type SafeMap struct {
	lock sync.Mutex
	data map[string]interface{}
}

var m sync.Mutex
var Sm = &SafeMap{data: map[string]interface{}{}}

func main() {
	runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	//模拟两个处理请求的goroutine
	go worker1()
	go worker2()
	go worker3()
	go worker4()
	//让pprof服务运行起来
	http.ListenAndServe(":6060", nil)
}

func worker4() {
	for {
		fmt.Println("xiaoming")
		time.Sleep(500 * time.Millisecond)
	}
}
func worker3() {
	m.Lock()
	go func() {
		time.Sleep(time.Second)
		m.Unlock()
	}()
	m.Lock()
	go func() {
		m.Lock()
	}()
	time.Sleep(time.Second)
	m.Unlock()
}

func worker1() {
	for {
		Sm.lock.Lock()
		defer Sm.lock.Unlock()
		Sm.data["test"] = 1
		time.Sleep(10 * time.Second)
	}
}

func worker2() {
	for {
		Sm.lock.Lock()
		defer Sm.lock.Unlock()
		Sm.data["test"] = 2
		time.Sleep(10 * time.Second)
	}
}

var datas = []string{}

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}
