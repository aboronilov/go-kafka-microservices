package main

import (
	"github.com/aboronilov/go-kafka-microservices/aggregator/client"
	"github.com/sirupsen/logrus"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000"
)

func main() {
	calcService := NewCalculatorService()
	calcService = NewLogMiddleware(calcService)

	httpClient := client.NewClient(aggregatorEndpoint)

	// grpcClient, err := client.NewGrpcClient(aggregatorEndpoint)
	// if err != nil {
	// 	logrus.Fatalf("Error creating gRPC client: %v\n", err)
	// 	return
	// }

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calcService, httpClient)
	if err != nil {
		logrus.Fatalf("Error creating Kafka consumer: %v\n", err)
		return
	}

	kafkaConsumer.Start()
}
