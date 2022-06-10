package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/BurntSushi/toml"
)

var counter int
var mutex = &sync.RWMutex{}

type Config struct {
	Port string `toml:"port"`
}

func hello(w http.ResponseWriter, _ *http.Request) {
	log.Println("call /hello")
	fmt.Fprintf(w, "hello world!\n")
}

func incrementCounter(_ http.ResponseWriter, _ *http.Request) {
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func decrementCounter(_ http.ResponseWriter, _ *http.Request) {
	mutex.Lock()
	counter--
	mutex.Unlock()
}

func getCounter(w http.ResponseWriter, _ *http.Request) {
	mutex.RLock()
	fmt.Fprintf(w, "counter = %v\n", counter)
	mutex.RUnlock()
}

func main() {
	cfg := &Config{}
	if _, err := toml.DecodeFile("config.toml", cfg); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", hello)

	http.HandleFunc("/inc", incrementCounter)

	http.HandleFunc("/dec", decrementCounter)

	http.HandleFunc("/count", getCounter)

	log.Printf("start http server on port = %v\n", cfg.Port)

	log.Fatal(http.ListenAndServe(cfg.Port, nil))

}
