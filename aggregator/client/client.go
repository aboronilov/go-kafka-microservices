package client

import (
	"context"

	"github.com/aboronilov/go-kafka-microservices/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
