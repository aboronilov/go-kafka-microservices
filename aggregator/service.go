package main

import (
	"fmt"

	"github.com/aboronilov/go-kafka-microservices/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{store: store}
}

func (i *InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	fmt.Println("proceesing and inserting distance to storage:", dist)
	return i.store.Insert(dist)
}

func (i *InvoiceAggregator) CalculateInvoice(obu int) (*types.Invoice, error) {
	distance, err := i.store.Get(obu)
	if err != nil {
		return nil, err
	}

	return &types.Invoice{
		OBUID:         obu,
		TotalDistance: distance,
		TotalAmount:   distance * basePrice,
	}, nil
}
