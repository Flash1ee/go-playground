package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

const DONE = "done"

type bucket struct {
	msg   string
	value int64
}
type Channel struct {
	counter int64
	ch      chan *bucket
	incDone chan interface{}
}

func NewChannel() *Channel {
	c := &Channel{}
	c.ch = make(chan *bucket)
	c.incDone = make(chan interface{})
	return c
}
func (s *Channel) Inc() {
	s.ch <- &bucket{value: 1}
}

func (s *Channel) Sum() int64 {
	for val := range s.ch {
		if val.msg == DONE {
			return s.counter
		}
		s.counter += val.value
	}

	return s.counter
}

func ChannelWorker(job Channel, n int) int64 {
	go job.Sum()

	var wgWorkers sync.WaitGroup
	countWorkers := runtime.GOMAXPROCS(0)

	var curWorkers int64
	var cnt int64

	for cnt < int64(n) {
		if curWorkers >= int64(countWorkers) {
			runtime.Gosched()
		}

		curWorkers++
		cnt++
		wgWorkers.Add(1)
		go func() {
			defer func() {
				wgWorkers.Done()
				atomic.AddInt64(&curWorkers, -1)

			}()
			job.Inc()
		}()
	}

	//var wg sync.WaitGroup
	//wg.Add(n)
	//for i := 0; i < n; i++ {
	//	go func() {
	//		defer wg.Done()
	//		job.Inc()
	//	}()
	//}
	//wg.Wait()
	//fmt.Println(runtime.NumGoroutine())

	wgWorkers.Wait()
	job.ch <- &bucket{msg: DONE}
	return job.counter
}
