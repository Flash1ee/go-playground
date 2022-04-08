package workerpool

import (
	"runtime"
	"sync/atomic"
)

type WorkerPoolAtomic struct {
	result  []Result
	resChan chan Result

	aggregateDone bool
	countWorkers  int32
	limit         int
}

func NewPoolAtomic(limit int) *WorkerPoolAtomic {
	if limit <= 0 {
		limit = DefaultLimit
	}
	return &WorkerPoolAtomic{
		limit:   limit,
		resChan: make(chan Result),
	}
}

func (wp *WorkerPoolAtomic) GetNumInProgress() int32 {
	return atomic.LoadInt32(&wp.countWorkers)
}

func (wp *WorkerPoolAtomic) Stop() []Result {
	for wp.countWorkers != 0 {
		runtime.Gosched()
	}

	close(wp.resChan)

	for wp.aggregateDone != true {
		runtime.Gosched()
	}

	return wp.result
}

func (wp *WorkerPoolAtomic) Run(tasks []Task) {
	go wp.Aggregate()
	for _, task := range tasks {
		wp.run(task)
	}
}

func (wp *WorkerPoolAtomic) run(task Task) {
	for {
		if wp.GetNumInProgress() < int32(wp.limit) {
			atomic.AddInt32(&wp.countWorkers, 1)

			go func(t Task) {
				defer func() {
					atomic.AddInt32(&wp.countWorkers, -1)
				}()

				task.Work(wp.resChan)
			}(task)
			break
		}
		runtime.Gosched()
	}

}
func (wp *WorkerPoolAtomic) Aggregate() {
	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	wp.aggregateDone = true
}
