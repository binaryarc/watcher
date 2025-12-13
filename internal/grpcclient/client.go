package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/binaryarc/watcher/internal/auth"
	"github.com/binaryarc/watcher/internal/detector"
	pb "github.com/binaryarc/watcher/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps gRPC client connection
type Client struct {
	conn   *grpc.ClientConn
	client pb.WatcherServiceClient
	apiKey string
}

// NewClient creates a new gRPC client
func NewClient(host string, apiKey string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewWatcherServiceClient(conn),
		apiKey: apiKey,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ObserveRuntimes fetches runtime information from remote server
func (c *Client) ObserveRuntimes(ctx context.Context) ([]*detector.Runtime, error) {
	if c.apiKey != "" {
		ctx = auth.InjectAPIKey(ctx, c.apiKey)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &pb.ObserveRequest{}
	resp, err := c.client.ObserveRuntimes(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	runtimes := make([]*detector.Runtime, 0, len(resp.Runtimes))
	for _, protoRuntime := range resp.Runtimes {
		if protoRuntime.Found {
			runtimes = append(runtimes, &detector.Runtime{
				Name:    protoRuntime.Name,
				Version: protoRuntime.Version,
				Path:    protoRuntime.Path,
				Found:   protoRuntime.Found,
			})
		}
	}

	return runtimes, nil
}

// ObserveRuntime fetches specific runtime information from remote server
func (c *Client) ObserveRuntime(ctx context.Context, name string) (*detector.Runtime, error) {
	runtimes, err := c.ObserveRuntimes(ctx)
	if err != nil {
		return nil, err
	}

	for _, runtime := range runtimes {
		if runtime.Name == name {
			return runtime, nil
		}
	}

	return &detector.Runtime{
		Name:  name,
		Found: false,
	}, nil
}
