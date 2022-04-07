package main

import (
	"fmt"
	"log"
	"time"

	"workerpool/workerpool"
)

func main() {
	tasks := []workerpool.Task{
		workerpool.Site{URL: "https://www.google.com"},
		workerpool.Site{URL: "https://facebook.com"},
		workerpool.Site{URL: "https://www.yandex.ru"},
		workerpool.Site{URL: "https://www.mail.ru"},
	}
	for i := 0; i < 4; i++ {
		tasks = append(tasks, tasks...)
	}

	//wp := workerpool.NewSyncPool(tasks, workerpool.MaxGoroutines)
	//wp.Run()

	//wpAtomic := workerpool.NewPoolAtomic(tasks, workerpool.MaxGoroutines)
	//wpAtomic.Run()

	fmt.Println("Count tasks ", len(tasks))

	//res := wpAtomic.Stop()
	//res := wp.Stop()
	wp := workerpool.NewPoolChannel(128)
	start := time.Now()
	if _, err := wp.Run(tasks); err != nil {
		log.Fatalln(err)
	}
	res := wp.Stop()

	fmt.Println("Count processing tasks ", len(res))
	fmt.Printf("Executing time = %v\n", time.Now().Sub(start).Seconds())
}
