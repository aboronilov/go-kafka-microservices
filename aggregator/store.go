package main

import "github.com/aboronilov/go-kafka-microservices/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(dist types.Distance) error {
	return nil
}
