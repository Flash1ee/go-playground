package workerpool

import (
	"errors"
	"log"
	"runtime"
	"sync/atomic"
)

type WorkerPoolChannel struct {
	limit       int
	workers     chan int
	cntExecutes int32

	resChan chan Result
	resDone bool
	result  []Result
}

func NewPoolChannel(limit int) *WorkerPoolChannel {
	if limit <= 0 {
		limit = DefaultLimit
	}
	pool := &WorkerPoolChannel{
		limit:       limit,
		workers:     make(chan int, limit),
		cntExecutes: 0,
		resChan:     make(chan Result),
	}

	for i := 0; i < limit; i++ {
		pool.workers <- i
	}

	return pool
}
func (wp *WorkerPoolChannel) GetNumInProgress() int32 {
	return atomic.LoadInt32(&wp.cntExecutes)
}

func (wp *WorkerPoolChannel) Stop() []Result {
	if err := wp.waitAndClose(); err != nil {
		log.Println(err)
	}
	wp.aggregateDone()

	return wp.result
}

func (wp *WorkerPoolChannel) Run(tasks []Task) (int, error) {
	go wp.aggregate()

	for _, task := range tasks {
		_, err := wp.run(task)
		//log.Printf("worker with ID = %v completed\n", ID)
		if err != nil {
			wp.Stop()
			return -1, err
		}
	}
	return 0, nil
}
func (wp *WorkerPoolChannel) run(task Task) (int, error) {
	wrkID, available := <-wp.workers
	if !available {
		return -1, errors.New("can not get worker from channel")
	}
	atomic.AddInt32(&wp.cntExecutes, 1)

	go func(t Task) {
		defer func() {
			atomic.AddInt32(&wp.cntExecutes, -1)
			wp.workers <- wrkID
		}()

		task.Work(wp.resChan)
	}(task)

	return wrkID, nil
}

func (wp *WorkerPoolChannel) aggregate() {
	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	wp.resDone = true
}
func (wp *WorkerPoolChannel) aggregateDone() {
	close(wp.resChan)
	for !wp.resDone {
		runtime.Gosched()
	}
}

func (wp *WorkerPoolChannel) waitAndClose() error {
	for i := 0; i < wp.limit; i++ {
		<-wp.workers
	}
	return wp.close()
}

func (wp *WorkerPoolChannel) close() error {
	close(wp.workers)
	return nil
}
