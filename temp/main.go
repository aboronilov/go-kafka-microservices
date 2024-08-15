package main

import (
	"context"
	"log"
	"time"

	"github.com/aboronilov/go-kafka-microservices/aggregator/client"
	"github.com/aboronilov/go-kafka-microservices/types"
)

func main() {
	c, err := client.NewGrpcClient(":3001")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
