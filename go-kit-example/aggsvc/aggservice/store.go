package aggservice

import (
	"fmt"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{
		"ObuID":     dist.OBUID,
		"Distance":  dist.Value,
		"Timestamp": dist.Unix,
	})
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

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(dist types.Distance) error {
	m.data[dist.OBUID] += dist.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("invoice distance not found for OBU %d", id)
	}

	return dist, nil
}
