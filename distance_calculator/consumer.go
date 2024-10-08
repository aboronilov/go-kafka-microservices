package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aboronilov/go-kafka-microservices/aggregator/client"
	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, client client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   client,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka consumer started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) Stop() {
	logrus.Info("kafka consumer stopping")
	c.isRunning = false
	c.consumer.Close()
	time.Sleep(time.Second) // wait for goroutine to finish before exiting the main function
	logrus.Info("kafka consumer stopped")
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error: %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("unmarshal error: %s", err)
			logrus.WithFields(logrus.Fields{
				"error":     err.Error(),
				"requsetID": data.RequestID,
			})
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("error calculating distance: %s", err)
			continue
		}

		dist := &types.AggregateRequest{
			Value: distance,
			ObuID: int32(data.OBUID),
			Unix:  time.Now().Unix(),
			// RequestID: data.RequestID,
		}

		if err := c.aggClient.Aggregate(context.Background(), dist); err != nil {
			logrus.Errorf("error aggregating invoice: %s", err)
			continue
		}
	}
}
