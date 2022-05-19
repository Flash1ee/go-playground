package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var counter int
var mutex = &sync.RWMutex{}

var port = ":8088"

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!\n")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func decrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter--
	mutex.Unlock()
}

func getCounter(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	fmt.Fprintf(w, "counter = %v\n", counter)
	mutex.RUnlock()
}

func main() {
	http.HandleFunc("/", hello)

	http.HandleFunc("/inc", incrementCounter)

	http.HandleFunc("/dec", decrementCounter)

	http.HandleFunc("/count", getCounter)

	log.Printf("start http server on port = %v\n", port)

	log.Fatal(http.ListenAndServe(port, nil))

}
