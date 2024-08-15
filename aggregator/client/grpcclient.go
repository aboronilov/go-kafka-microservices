package client

import (
	"github.com/aboronilov/go-kafka-microservices/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGrpcClient(endpoint string) (*GRPCClient, error) {
	// conn, err := grpc.Dial(endpoint)
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}, nil
}
