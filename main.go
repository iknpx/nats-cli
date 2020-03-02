package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"sync"
)

var usage = "Usage \n\n" +
	"emit:             nats [channel] [message]\n" +
	"subscribe:        nats [channel]\n"

func handlePanic() string {
	fmt.Print(usage)
	os.Exit(0)

	return ""
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getActionType() string {
	if len(os.Args) == 2 {
		return "s"
	}

	if len(os.Args) == 3 {
		return "e"
	}

	return handlePanic()
}

func subscribe() {
	nc, connectionError := nats.Connect(nats.DefaultURL)
	handleError(connectionError)

	wg := sync.WaitGroup{}
	wg.Add(1)

	if _, err := nc.Subscribe(os.Args[1], func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		wg.Done()
	}); err != nil {
		handleError(err)
	}

	wg.Wait()
}

func emit() {
	nc, connectionError := nats.Connect(nats.DefaultURL)
	handleError(connectionError)

	defer nc.Close()

	if err := nc.Publish(os.Args[1], []byte(os.Args[2])); err != nil {
		handleError(err)
	}
}

func main() {
	action := getActionType()

	switch action {
	case "s":
		subscribe()
		break

	case "e":
		emit()
		break
	}
}
