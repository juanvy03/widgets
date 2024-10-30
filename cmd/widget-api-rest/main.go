package main

import (
	"log"
	"spike/internal/redisclient"
	"spike/internal/service"
)

func main() {
	log.Println("API server is UP")

	stop := make(chan bool)

	go func() {
		service.InitHttpServer()
	}()

	/*
	 In the best scenario would be good to be another microservice with it's own scaling.
	*/
	go func() {
		redisclient.InitRedisClient().StreamConsumer("events")
	}()
	<-stop
}
