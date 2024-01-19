package main

import (
	"fmt"
	"tolling/types"

)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalcualteInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	fmt.Println("processing and inserting distance in the storage: ", dist)
	return i.store.Insert(dist)
}

func (i *InvoiceAggregator) CalcualteInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("obu id (%d) not found", id)
	}
	invoice := &types.Invoice{
		OBUID: id,
		TotalDistance: dist,
		TotalAmount: basePrice * dist,
	}
	return invoice, nil
}
