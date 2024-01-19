package main

import (
	"context"
	"tolling/types"
)

type GrpcAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGrpcServer(svc Aggregator) *GrpcAggregatorServer {
	return &GrpcAggregatorServer{
		svc: svc,
	}
}

func (s *GrpcAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
