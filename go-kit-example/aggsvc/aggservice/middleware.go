package aggservice

import (
	"context"
	"fmt"

	"github.com/aboronilov/go-kafka-microservices/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{next: next}
	}
}

func (lm *loggingMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	fmt.Println("Processing and inserting distance to storage:", dist)
	return nil
}

func (lm *loggingMiddleware) Calculate(_ context.Context, obu int) (*types.Invoice, error) {
	return nil, fmt.Errorf("invoice calculation not implemented")
}

type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}

func (im *instrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	fmt.Println("Processing and inserting distance to storage:", dist)
	return nil
}

func (im *instrumentationMiddleware) Calculate(_ context.Context, obu int) (*types.Invoice, error) {
	return nil, fmt.Errorf("invoice calculation not implemented")
}
