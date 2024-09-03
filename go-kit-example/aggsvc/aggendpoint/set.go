package aggendpoint

import (
	"context"

	aggserrvice "github.com/aboronilov/go-kafka-microservices/go-kit-example/aggsvc/aggservice"
	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s *Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		Value: dist.Value,
		OBUID: dist.OBUID,
		Unix:  dist.Unix,
	})

	return err
}

func (s *Set) Calculate(ctx context.Context, obu int) (*types.Invoice, error) {
	response, err := s.CalculateEndpoint(ctx, CalculateRequest{OBUID: obu})
	if err != nil {
		return nil, err
	}

	result := response.(CalculateResponse)
	if result.Err != nil {
		return nil, result.Err
	}

	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func MakeAggregateEndpoint(s aggserrvice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			Value: req.Value,
			OBUID: req.OBUID,
			Unix:  req.Unix,
		})
		return AggregateResponse{
			Err: err,
		}, err
	}
}

func MakeCalculateEndpoint(s aggserrvice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		v, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			OBUID:         req.OBUID,
			TotalDistance: v.TotalDistance,
			TotalAmount:   v.TotalAmount,
			Err:           err,
		}, err
	}
}
