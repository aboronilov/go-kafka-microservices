package aggservice

import (
	"context"

	"github.com/aboronilov/go-kafka-microservices/types"
)

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store Storer
}

func (s BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return s.store.Insert(dist)
}

func (s BasicService) Calculate(_ context.Context, obu int) (*types.Invoice, error) {
	dist, err := s.store.Get(obu)
	if err != nil {
		return nil, err
	}

	return &types.Invoice{
		OBUID:         obu,
		TotalDistance: dist,
		TotalAmount:   dist * basePrice,
	}, nil
}

func newBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func NewAggregatorSeervice() Service {
	var svc Service
	{
		svc = newBasicService(NewMemoryStore())
		svc = newLoggingMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}

	return svc
}
