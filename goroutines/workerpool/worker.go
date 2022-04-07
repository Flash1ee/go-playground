package workerpool

import "sync"

type worker struct {
	resChan chan Result
}

func NewWorker(resChan chan Result) *worker {
	return &worker{
		resChan: resChan,
	}
}
func (wrk *worker) Start(src <-chan Task, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range src {
			task.Work(wrk.resChan)
		}
	}()
}
