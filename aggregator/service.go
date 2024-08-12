package main

import (
	"fmt"

	"github.com/aboronilov/go-kafka-microservices/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{store: store}
}

func (i *InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	fmt.Println("proceesing and inserting distance to storage:", dist)
	return i.store.Insert(dist)
}
