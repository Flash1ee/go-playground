package workerpool

import (
	"log"
	"sync/atomic"
)

type WorkerPoolChannel struct {
	limit       int
	workers     chan struct{}
	cntExecutes int32

	resChan chan Result
	resDone chan int
	result  []Result
}

func NewPoolChannel(limit int) *WorkerPoolChannel {
	if limit <= 0 {
		limit = DefaultLimit
	}
	pool := &WorkerPoolChannel{
		limit:       limit,
		workers:     make(chan struct{}, limit),
		cntExecutes: 0,
		resChan:     make(chan Result),
		resDone:     make(chan int),
	}

	//for i := 0; i < limit; i++ {
	//	pool.workers <- i
	//}

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
		err := wp.run(task)
		//log.Printf("worker with ID = %v completed\n", ID)
		if err != nil {
			wp.Stop()
			return -1, err
		}
	}
	return 0, nil
}
func (wp *WorkerPoolChannel) run(task Task) error {
	wp.workers <- struct{}{}
	//wrkID, available := <-wp.workers
	//if !available {
	//	return -1, errors.New("can not get worker from channel")
	//}
	atomic.AddInt32(&wp.cntExecutes, 1)

	go func(t Task) {
		defer func() {
			atomic.AddInt32(&wp.cntExecutes, -1)
			//wp.workers <- wrkID
			<-wp.workers
		}()

		task.Work(wp.resChan)
	}(task)

	return nil
}

func (wp *WorkerPoolChannel) aggregate() {
	for value := range wp.resChan {
		wp.result = append(wp.result, value)
	}
	close(wp.resDone)
}
func (wp *WorkerPoolChannel) aggregateDone() {
	close(wp.resChan)
	<-wp.resDone

}

func (wp *WorkerPoolChannel) waitAndClose() error {
	for i := 0; i < len(wp.workers); i++ {
		wp.workers <- struct{}{}
	}
	return wp.close()
}

func (wp *WorkerPoolChannel) close() error {
	close(wp.workers)
	return nil
}
