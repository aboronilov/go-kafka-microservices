package main

import (
	"fmt"

	"github.com/aboronilov/go-kafka-microservices/types"
)

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
		return 0.0, fmt.Errorf("invoice distance not found for OBU", id)
	}

	return dist, nil
}
