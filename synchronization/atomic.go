package main

import "sync/atomic"

type Atomic struct {
	counter int64
}

func NewAtomic() *Atomic {
	return &Atomic{}
}

func (s *Atomic) Inc() {
	atomic.AddInt64(&s.counter, 1)
}

func (s *Atomic) Sum() int64 {
	return atomic.LoadInt64(&s.counter)
}
