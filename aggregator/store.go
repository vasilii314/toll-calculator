package main

import (
	"fmt"
	"tolling/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(obuId int) (float64, error) {
	dist, ok := m.data[obuId]
	if !ok {
		return 0.0, fmt.Errorf("distance not found for id %d", obuId)
	}
	return dist, nil
}