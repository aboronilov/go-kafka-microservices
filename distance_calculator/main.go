package main

import "github.com/sirupsen/logrus"

const kafkaTopic = "obudata"

func main() {
	calcService := NewCalculatorService()
	calcService = NewLogMiddleware(calcService)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calcService)
	if err != nil {
		logrus.Fatalf("Error creating Kafka consumer: %v\n", err)
		return
	}

	kafkaConsumer.Start()
}
