package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

type Config struct {
	Port string `toml:"bind_addr"`
}

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/cfg.toml", "path to config file")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!\n")
}

var counter int
var mutex = &sync.RWMutex{}

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
	flag.Parse()

	config := Config{}

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		log.Fatal(err, configPath)
	}

	http.HandleFunc("/", hello)

	http.HandleFunc("/inc", incrementCounter)

	http.HandleFunc("/dec", decrementCounter)

	http.HandleFunc("/count", getCounter)

	log.Printf("start http server on port = %v\n", config.Port)

	log.Fatal(http.ListenAndServe(config.Port, nil))

}
