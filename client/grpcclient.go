package client

import (
	"context"
	"tolling/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	clent    types.AggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GrpcClient{
		Endpoint: endpoint,
		clent:    c,
	}, nil
}

func (c *GrpcClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.clent.Aggregate(ctx, req)
	return err
}

func (c *GrpcClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID: id,
		TotalDistance: 12432,
		TotalAmount: 124.234,
	}, nil
}
