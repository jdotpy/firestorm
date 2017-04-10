package main

import "net/http"
import "fmt"
import "time"

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

type HttpRequest struct {
	url      string
}

func fetcher(requests chan *HttpRequest, responses chan *HttpResponse) {
	time.Sleep(5 * time.Second)
	for {
		request := <-requests
		fmt.Println("GET " + request.url)
		resp, err := http.Get(request.url)
		resp.Body.Close()
		responses <-&HttpResponse{request.url, resp, err}
	}
}

func RapidFire(workerCount int, bufferSize int) (chan *HttpRequest, chan *HttpResponse) {
	requests := make(chan *HttpRequest, bufferSize)
	responses := make(chan *HttpResponse, bufferSize)

    fmt.Println("Starting workers...")
	for i := 0; i < workerCount; i++ {
		fmt.Printf("Starting worker #%d\n", i + 1)
		go fetcher(requests, responses)
	}

    fmt.Println("Ready for requests")
	return requests, responses
}
