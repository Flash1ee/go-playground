package workerpool

import (
	"runtime"
)

type WorkerPoolAtomic struct {
	tasks  []Task
	result []Result

	resChan       chan Result
	aggregateDone bool
	maxWorkers    int
	countWorkers  int64
}

func NewPoolAtomic(tasks []Task, maxWorkers int) *WorkerPoolAtomic {
	return &WorkerPoolAtomic{
		tasks:      tasks,
		maxWorkers: maxWorkers,
		resChan:    make(chan Result),
	}
}
func (wp *WorkerPoolAtomic) Stop() []Result {
	if wp.countWorkers != 0 {
		runtime.Gosched()
	}
	close(wp.resChan)
	if wp.aggregateDone != true {
		runtime.Gosched()
	}

	return wp.result

}

func (wp *WorkerPoolAtomic) Run() {
	dataChan := make(chan Task)

	go wp.Aggregate()

	for i := 0; i < wp.maxWorkers; i++ {
		wrk := NewWorker(wp.resChan)
		wrk.Start(dataChan, nil)
	}

	for _, val := range wp.tasks {
		dataChan <- val
	}
	close(dataChan)
}
func (wp *WorkerPoolAtomic) Aggregate() {
	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	wp.aggregateDone = true
}
