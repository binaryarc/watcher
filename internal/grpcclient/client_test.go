package grpcclient

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/binaryarc/watcher/internal/auth"
	pb "github.com/binaryarc/watcher/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeWatcherServer struct {
	pb.UnimplementedWatcherServiceServer
	resp      *pb.ObserveResponse
	lastKey   string
	callCount int
}

func (f *fakeWatcherServer) ObserveRuntimes(ctx context.Context, req *pb.ObserveRequest) (*pb.ObserveResponse, error) {
	f.callCount++
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(auth.APIKeyHeader)
		if len(values) > 0 {
			f.lastKey = values[0]
		}
	}
	return f.resp, nil
}

func startFakeServer(t *testing.T, resp *pb.ObserveResponse) (addr string, cleanup func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	fake := &fakeWatcherServer{resp: resp}
	pb.RegisterWatcherServiceServer(grpcServer, fake)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	cleanup = func() {
		grpcServer.Stop()
		_ = lis.Close()
	}

	return lis.Addr().String(), cleanup
}

func TestClientObserveRuntimes(t *testing.T) {
	t.Parallel()

	resp := &pb.ObserveResponse{
		Runtimes: []*pb.Runtime{
			{Name: "python", Version: "3.10.12", Path: "/usr/bin/python3", Found: true},
		},
	}
	addr, cleanup := startFakeServer(t, resp)
	defer cleanup()

	client, err := NewClient(addr, "secret-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	runtimes, err := client.ObserveRuntimes(ctx)
	if err != nil {
		t.Fatalf("ObserveRuntimes() error = %v", err)
	}

	if len(runtimes) != 1 {
		t.Fatalf("expected 1 runtime, got %d", len(runtimes))
	}

	rt := runtimes[0]
	if rt.Name != "python" || rt.Version != "3.10.12" || rt.Path != "/usr/bin/python3" {
		t.Fatalf("unexpected runtime: %+v", rt)
	}
}

func TestClientObserveRuntimeNotFound(t *testing.T) {
	t.Parallel()

	resp := &pb.ObserveResponse{
		Runtimes: []*pb.Runtime{
			{Name: "java", Version: "17", Path: "/opt/java/bin/java", Found: true},
		},
	}
	addr, cleanup := startFakeServer(t, resp)
	defer cleanup()

	client, err := NewClient(addr, "")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	java, err := client.ObserveRuntime(ctx, "java")
	if err != nil {
		t.Fatalf("ObserveRuntime(java) error = %v", err)
	}
	if !java.Found {
		t.Fatalf("expected java runtime to be found")
	}

	goRuntime, err := client.ObserveRuntime(ctx, "go")
	if err != nil {
		t.Fatalf("ObserveRuntime(go) error = %v", err)
	}
	if goRuntime.Found {
		t.Fatalf("expected go runtime to be missing")
	}
}

func TestClientInjectsAPIKey(t *testing.T) {
	t.Parallel()

	resp := &pb.ObserveResponse{}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	fake := &fakeWatcherServer{resp: resp}
	pb.RegisterWatcherServiceServer(grpcServer, fake)

	go func() {
		_ = grpcServer.Serve(lis)
	}()
	defer func() {
		grpcServer.Stop()
		_ = lis.Close()
	}()

	client, err := NewClient(lis.Addr().String(), "api-key-value")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := client.ObserveRuntimes(ctx); err != nil {
		t.Fatalf("ObserveRuntimes() error = %v", err)
	}

	if fake.lastKey != "api-key-value" {
		t.Fatalf("expected API key header to be forwarded, got %q", fake.lastKey)
	}
	if fake.callCount == 0 {
		t.Fatalf("fake server was not called")
	}
}
