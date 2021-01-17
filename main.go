package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/willdot/brabbit/rabbit"
	"github.com/willdot/brabbit/service"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "an array of values in a flag"
}

func (i *arrayFlags) Set(value string) error {
	values := strings.Split(value, ",")
	for _, v := range values {
		if v == "" {
			continue
		}

		*i = append(*i, v)
	}
	return nil
}

func main() {
	queue := flag.String("queue", "", "t'he queue to send the message to")
	bodyFileName := flag.String("body", "", "the file name of the body to send in the message")
	headersFileName := flag.String("headers", "", "the file name of the header to send in the message (optional)")
	repeat := flag.Int("repeat", 1, "the number of times to send the message")

	var bodyFileNames arrayFlags
	flag.Var(&bodyFileNames, "bodies", "a comma seperated list of filenames of bodies to send as messages. (This overrides the 'body' flag)")

	flag.Parse()

	if *queue == "" {
		log.Fatalf("please provide a queue to send the message to")
	}

	if len(bodyFileNames) == 0 && *bodyFileName == "" {
		log.Fatalf("please provide either the 'body' or 'bodies' flags")
	}

	// if the bodies flag wasn't used, add the body flag to use. Otherwise the body flag will be ignored and the bodies flag will be used
	if len(bodyFileNames) == 0 {
		bodyFileNames = append(bodyFileNames, *bodyFileName)
	}

	if err := run(*queue, *headersFileName, bodyFileNames, *repeat); err != nil {
		log.Fatalf("error running: %s", err.Error())
	}
}

func run(queue, headersFileName string, bodyFilenames []string, repeat int) error {
	publisher := rabbit.NewRabbitPublisher()
	defer publisher.Shutdown()

	srv := service.NewService(publisher)

	bodies, err := getAllBodies(bodyFilenames)
	if err != nil {
		return err
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

	for _, body := range bodies {
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
	}

	return nil
}

func getAllBodies(filename []string) ([][]byte, error) {
	var result [][]byte

	for _, bodyPath := range filename {
		body, err := ioutil.ReadFile(bodyPath)
		if err != nil {
			return nil, fmt.Errorf("error getting body file: %s", err.Error())
		}

		result = append(result, body)
	}

	return result, nil
}
