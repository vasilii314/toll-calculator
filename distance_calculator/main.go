package main

import (
	"log"
	"tolling/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:5050/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	grpcClient, err := client.NewGrpcClient(":5051")
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
