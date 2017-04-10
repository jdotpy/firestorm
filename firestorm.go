package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"

	//"gopkg.in/urfave/cli.v2"
	//"gopkg.in/yaml.v2"
)

func requestEmitter(requests chan *HttpRequest) {
	url := "http://localhost:8000/"
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Queue", url)
		requests <-&HttpRequest{url}
		fmt.Println("Queue", url)
		requests <-&HttpRequest{url}
		fmt.Println("Queue", url)
		requests <-&HttpRequest{url}
		fmt.Println("Queue", url)
		requests <-&HttpRequest{url}
		fmt.Println("Queue", url)
		requests <-&HttpRequest{url}
	}
}

func fire() {
	// Prep the signals channel to listen for keyboard interrupts
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Make the channels by which we communicate with the http fetcher
	requests, responses := RapidFire(5, 100)

	// Now start the emitter
    go requestEmitter(requests)

	// Check our channel for responses and save them to results
	results := []*HttpResponse{}
	for {
		select {
			case sig := <-signals:
				fmt.Println("received signal", sig)
				os.Exit(1)
			case response := <-responses:
				fmt.Println("got response")
				results = append(results, response)
		}
	}
    fmt.Println("done")
}

func main() {
	//cmd := &cli.App{
	//	Name: "Firestorm",
	//	Usage: "Generate HTTP requests",
	//	Action: fire,
	//}
	//cmd.Run(os.Args)
	fire()
}

