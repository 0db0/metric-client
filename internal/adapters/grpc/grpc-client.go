package grpc

import (
	"context"
	pb "github.com/0db0/metric-server/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"metric-client/config"
	"metric-client/internal/models"
)

type Client struct {
	pbc pb.MetricClient
}

func NewClient(ctx context.Context, cfg config.Config) *Client {
	conn, err := grpc.DialContext(ctx, cfg.GRPCClient.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &Client{
		pbc: pb.NewMetricClient(conn),
	}
}

func (c Client) SendMetric(ctx context.Context, ch <-chan []models.Metric) {

	for {
		select {
		case <-ctx.Done():
			return
		case metrics, ok := <-ch:
			if !ok {
				return
			}

			go func(metrics []models.Metric) {
				request := c.createRequest(metrics)
				_, err := c.pbc.BatchCollect(ctx, request)

				if err != nil {
					// log
				}
			}(metrics)

		}
	}
}

func (c Client) createRequest(metrics []models.Metric) *pb.BatchCollectMetricRequest {
	var metricRequestList []*pb.CollectMetricRequest

	for _, metric := range metrics {
		collectMetricRequest := &pb.CollectMetricRequest{
			Name:  metric.ID,
			Type:  metric.MType,
			Delta: c.prepareDelta(metric.Delta),
			Value: c.prepareValue(metric.Value),
		}

		metricRequestList = append(metricRequestList, collectMetricRequest)
	}

	return &pb.BatchCollectMetricRequest{
		Metrics: metricRequestList,
	}
}

func (c Client) prepareDelta(delta *int64) *wrapperspb.Int64Value {
	if delta == nil {
		return nil
	}

	return &wrapperspb.Int64Value{Value: *delta}
}

func (c Client) prepareValue(value *float64) *wrapperspb.DoubleValue {
	if value == nil {
		return nil
	}

	return &wrapperspb.DoubleValue{Value: *value}
}
