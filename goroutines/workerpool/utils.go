package workerpool

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
)

type Task interface {
	Work(dst chan<- Result)
}

type Worker interface {
	Start(src <-chan Site)
}

type Site struct {
	URL string
}
type Result struct {
	Status int
	Body   string
}

func (s Site) Work(dst chan<- Result) {
	client := http.Client{
		Timeout: timeoutTime,
	}

	res, err := client.Get(s.URL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	body := &bytes.Buffer{}
	reader := bufio.NewReader(res.Body)
	io.Copy(body, reader)

	dst <- Result{
		Status: res.StatusCode,
		Body:   body.String(),
	}
}
