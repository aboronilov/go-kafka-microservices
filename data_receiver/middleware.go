package main

import (
	"time"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"OBUID": data.OBUID,
			"Lat":   data.Lat,
			"Long":  data.Long,
			"Took":  time.Since(start),
		}).Info("producing to Kafka")
	}(time.Now())

	return m.next.ProduceData(data)
}
