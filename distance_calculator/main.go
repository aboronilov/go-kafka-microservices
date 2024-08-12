package main

import (
	"github.com/aboronilov/go-kafka-microservices/aggregator/client"
	"github.com/sirupsen/logrus"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000/aggregate"
)

func main() {
	calcService := NewCalculatorService()
	calcService = NewLogMiddleware(calcService)
	aggClient := client.NewClient(aggregatorEndpoint)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calcService, aggClient)
	if err != nil {
		logrus.Fatalf("Error creating Kafka consumer: %v\n", err)
		return
	}

	kafkaConsumer.Start()
}
