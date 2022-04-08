package main

import (
	"fmt"
	"time"

	"workerpool/workerpool"
)

func main() {
	tasks := []workerpool.Task{
		workerpool.Site{URL: "https://www.google.com"},
		workerpool.Site{URL: "https://avito.ru"},
		workerpool.Site{URL: "https://www.yandex.ru"},
		workerpool.Site{URL: "https://www.mail.ru"},
	}
	for i := 0; i < 4; i++ {
		tasks = append(tasks, tasks...)
	}

	fmt.Println("Count tasks ", len(tasks))
	start := time.Now()
	wp := workerpool.NewPoolAtomic(128)
	wp.Run(tasks)
	res := wp.Stop()

	fmt.Println("Count processing tasks ", len(res))
	fmt.Printf("Executing time = %v\n", time.Now().Sub(start).Seconds())
}
