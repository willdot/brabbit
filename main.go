package main

import (
	"log"

	"github.com/willdot/brabbit/rabbit"
	"github.com/willdot/brabbit/service"
)

func main() {
	publisher := rabbit.NewRabbitPublisher()
	defer publisher.Shutdown()

	srv := service.NewService(publisher)

	msg := []byte(`{"name":"will"}`)
	headers := map[string]interface{}{"one": "two"}

	err := srv.SendMessage(msg, headers, 1)
	if err != nil {
		log.Fatalf("error calling publish: %s", err.Error())
	}
}
