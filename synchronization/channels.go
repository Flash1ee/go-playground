package main

type Channel struct {
	counter int64
	ch      chan int64
	incDone chan int64
}

func NewChannel() *Channel {
	c := &Channel{}
	c.ch = make(chan int64)
	c.incDone = make(chan int64)

	go c.worker()

	return c
}
func (s *Channel) Inc() {
	s.ch <- 1
}

func (s *Channel) Sum() int64 {
	close(s.ch)

	return <-s.incDone
}

func (s *Channel) worker() {
	for val := range s.ch {
		s.counter += val
	}
	s.incDone <- s.counter
}
