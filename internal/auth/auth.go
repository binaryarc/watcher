package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const (
	APIKeyHeader = "x-api-key"
)

// Validator is an interface for validating API keys
type Validator interface {
	Validate(key string) bool
}

// ExtractAPIKey extracts API key from gRPC metadata
func ExtractAPIKey(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	values := md.Get(APIKeyHeader)
	if len(values) == 0 {
		return "", fmt.Errorf("no API key in metadata")
	}

	return values[0], nil
}

// InjectAPIKey adds API key to outgoing gRPC metadata
func InjectAPIKey(ctx context.Context, apiKey string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	md.Set(APIKeyHeader, apiKey)

	return metadata.NewOutgoingContext(ctx, md)
}
