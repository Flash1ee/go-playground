package workerpool

import (
	"sync"
)

type WorkerSyncPool struct {
	tasks  []Task
	result []Result

	resChan       chan Result
	aggregateDone chan struct{}
	maxWorkers    int
	wgWorkers     sync.WaitGroup
	wgAggregate   sync.WaitGroup
}

func NewSyncPool(tasks []Task, maxWorkers int) *WorkerSyncPool {
	return &WorkerSyncPool{
		tasks:      tasks,
		maxWorkers: maxWorkers,
		resChan:    make(chan Result),
	}
}
func (wp *WorkerSyncPool) Stop() []Result {
	wp.wgWorkers.Wait()
	close(wp.resChan)
	wp.wgAggregate.Wait()

	return wp.result

}

func (wp *WorkerSyncPool) Run() {
	dataChan := make(chan Task)

	go wp.Aggregate()

	for i := 0; i < wp.maxWorkers; i++ {
		wrk := NewWorker(wp.resChan)
		wrk.Start(dataChan, &wp.wgWorkers)
	}

	for _, val := range wp.tasks {
		dataChan <- val
	}
	close(dataChan)
}
func (wp *WorkerSyncPool) Aggregate() {
	wp.wgAggregate.Add(1)
	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	wp.wgAggregate.Done()
}
