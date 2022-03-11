package main

const buf = 4096

type BufChannel struct {
	counter int64
	ch      chan struct{}
	incDone chan int64
}

func NewBufChannel() *BufChannel {
	c := &BufChannel{}
	c.ch = make(chan struct{}, buf)
	c.incDone = make(chan int64)

	go c.worker()

	return c
}
func (s *BufChannel) Inc() {
	s.ch <- struct{}{}
}

func (s *BufChannel) Sum() int64 {
	close(s.ch)

	return <-s.incDone
}

func (s *BufChannel) worker() {
	for range s.ch {
		s.counter += 1
	}
	s.incDone <- s.counter
}
