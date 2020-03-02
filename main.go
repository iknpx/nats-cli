package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"sync"
)

var usage = "Usage \n\n" +
	"emit:             nats [channel] [message]\n" +
	"subscribe:        nats [channel]\n"

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	nc, connectionError := nats.Connect(nats.DefaultURL)
	handleError(connectionError)

	if len(os.Args) < 2 {
		log.Print(usage)
		os.Exit(0)
	}

	if len(os.Args) == 3 {
		defer nc.Close()

		if err := nc.Publish(os.Args[1], []byte(os.Args[2])); err != nil {
			handleError(err)
		}
	}

	if len(os.Args) == 2 {
		wg := sync.WaitGroup{}
		wg.Add(1)

		if _, err := nc.Subscribe(os.Args[1], func(m *nats.Msg) {
			log.Println(string(m.Data))
			wg.Done()
		}); err != nil {
			handleError(err)
		}

		wg.Wait()
	}
}
