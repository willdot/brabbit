package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/willdot/brabbit/rabbit"
	"github.com/willdot/brabbit/service"
)

func main() {
	queue := flag.String("queue", "", "t'he queue to send the message to")
	bodyFileName := flag.String("body", "", "the file name of the body to send in the message")
	headersFileName := flag.String("headers", "", "the file name of the header to send in the message (optional)")
	repeat := flag.Int("repeat", 1, "the number of times to send the message")
	flag.Parse()

	if *queue == "" {
		log.Fatalf("please provide a queue to send the message to")
	}

	if *bodyFileName == "" {
		log.Fatalf("please provide the body flag")
	}

	if err := run(*queue, *bodyFileName, *headersFileName, *repeat); err != nil {
		log.Fatalf("error running: %s", err.Error())
	}
}

func run(queue, bodyFileName, headersFileName string, repeat int) error {
	publisher := rabbit.NewRabbitPublisher()
	defer publisher.Shutdown()

	srv := service.NewService(publisher)

	body, err := ioutil.ReadFile(bodyFileName)
	if err != nil {
		return fmt.Errorf("error getting body file: %s", err.Error())
	}

	var headers map[string]interface{}

	if headersFileName != "" {
		headersBytes, err := ioutil.ReadFile(headersFileName)
		if err != nil {
			return fmt.Errorf("error getting headers file: %s", err.Error())
		}

		err = json.Unmarshal(headersBytes, &headers)
		if err != nil {
			return fmt.Errorf("error unmarshalling headers file: %s", err.Error())
		}
	}

	req := service.Request{
		Body:    body,
		Headers: headers,
		Repeat:  repeat,
		Queue:   queue,
	}

	err = srv.SendMessage(req)
	if err != nil {
		return fmt.Errorf("error calling publish: %s", err.Error())
	}

	return nil
}
