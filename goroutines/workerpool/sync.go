package workerpool

import (
	"runtime"
	"sync"
)

type WorkerPoolSync struct {
	result  []Result
	resChan chan Result

	wgAggregate sync.WaitGroup
	wgWorkers   sync.WaitGroup
	limit       int
}

func NewPoolSync(limit int) *WorkerPoolSync {
	if limit <= 0 {
		limit = DefaultLimit
	}
	return &WorkerPoolSync{
		limit:   limit,
		resChan: make(chan Result),
	}
}

func (wp *WorkerPoolSync) GetNumInProgress() int32 {
	return int32(runtime.NumGoroutine() - 2)
}

func (wp *WorkerPoolSync) Stop() []Result {
	wp.wgWorkers.Wait()

	close(wp.resChan)

	wp.wgAggregate.Wait()

	return wp.result
}

func (wp *WorkerPoolSync) Run(tasks []Task) {
	go wp.Aggregate()
	for _, task := range tasks {
		wp.run(task)
	}
}

func (wp *WorkerPoolSync) run(task Task) {
	for {
		if wp.GetNumInProgress() < int32(wp.limit) {
			wp.wgWorkers.Add(1)

			go func(t Task) {
				defer func() {
					wp.wgWorkers.Done()
				}()

				task.Work(wp.resChan)
			}(task)
			break
		}
		runtime.Gosched()
	}

}
func (wp *WorkerPoolSync) Aggregate() {
	wp.wgAggregate.Add(1)

	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	wp.wgAggregate.Done()
}
