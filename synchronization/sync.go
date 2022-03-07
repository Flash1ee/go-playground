package main

import "sync"

type BySync struct {
	counter int64
	mu      sync.Mutex
}

func NewSync() *BySync {
	return &BySync{}
}

func (s *BySync) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter += 1
}

func (s *BySync) Sum() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counter
}
