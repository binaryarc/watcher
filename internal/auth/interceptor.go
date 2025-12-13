package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC unary interceptor for API key validation
func UnaryServerInterceptor(validator Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		apiKey, err := ExtractAPIKey(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing API key")
		}

		if !validator.Validate(apiKey) {
			return nil, status.Error(codes.PermissionDenied, "invalid API key")
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream interceptor for API key validation
func StreamServerInterceptor(validator Validator) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		apiKey, err := ExtractAPIKey(ss.Context())
		if err != nil {
			return status.Error(codes.Unauthenticated, "missing API key")
		}

		if !validator.Validate(apiKey) {
			return status.Error(codes.PermissionDenied, "invalid API key")
		}

		return handler(srv, ss)
	}
}
