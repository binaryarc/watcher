package get

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/spf13/cobra"
)

func TestObserveLocalRuntimesFiltersFound(t *testing.T) {
	t.Parallel()

	originalDetectors := detectorsProvider
	defer func() { detectorsProvider = originalDetectors }()

	detectorsProvider = func() []detector.Detector {
		return []detector.Detector{
			&fakeDetector{name: "java", found: true, version: "17"},
			&fakeDetector{name: "missing", found: false},
		}
	}

	runtimes := observeLocalRuntimes("json")
	if len(runtimes) != 1 {
		t.Fatalf("expected 1 runtime, got %d", len(runtimes))
	}
	if runtimes[0].Name != "java" {
		t.Fatalf("unexpected runtime: %+v", runtimes[0])
	}
}

func TestRunGetRuntimesWithOutputFormats(t *testing.T) {
	t.Parallel()

	originalDetectors := detectorsProvider
	defer func() { detectorsProvider = originalDetectors }()

	detectorsProvider = func() []detector.Detector {
		return []detector.Detector{
			&fakeDetector{name: "java", found: true, version: "17"},
			&fakeDetector{name: "python", found: true, version: "3.10.12"},
		}
	}

	formats := []string{"json", "yaml", "table"}
	for _, format := range formats {
		buffer := &bytes.Buffer{}
		cmd := &cobra.Command{}
		cmd.SetOut(buffer)
		cmd.SetErr(buffer)
		cmd.Flags().String("output", format, "")
		cmd.Flags().String("host", "", "")

		runGetRuntimes(cmd, nil)
	}
}

func TestObserveRemoteRuntimesError(t *testing.T) {
	t.Parallel()

	originalFactory := grpcClientFactory
	defer func() { grpcClientFactory = originalFactory }()

	grpcClientFactory = func(host, apiKey string) (grpcRuntimeClient, error) {
		return nil, errors.New("connect failure")
	}

	if _, err := observeRemoteRuntimes("remote:9090", "", "table"); err == nil {
		t.Fatalf("expected error from observeRemoteRuntimes when client creation fails")
	}
}

func TestObserveRemoteRuntimesSuccess(t *testing.T) {
	t.Parallel()

	originalFactory := grpcClientFactory
	defer func() { grpcClientFactory = originalFactory }()

	expected := []*detector.Runtime{{Name: "java", Version: "17", Found: true}}
	mockClient := &fakeRuntimeClient{runtimes: expected}

	grpcClientFactory = func(host, apiKey string) (grpcRuntimeClient, error) {
		if host != "remote:9090" || apiKey != "secret" {
			t.Fatalf("unexpected host/apiKey: %s %s", host, apiKey)
		}
		return mockClient, nil
	}

	got, err := observeRemoteRuntimes("remote:9090", "secret", "json")
	if err != nil {
		t.Fatalf("observeRemoteRuntimes() error = %v", err)
	}
	if len(got) != len(expected) || got[0].Name != "java" {
		t.Fatalf("observeRemoteRuntimes() = %+v, want %+v", got, expected)
	}
	if !mockClient.closed {
		t.Fatalf("expected client.Close to be called")
	}
}

type fakeDetector struct {
	name    string
	version string
	path    string
	found   bool
}

func (f *fakeDetector) Name() string {
	return f.name
}

func (f *fakeDetector) Detect() (*detector.Runtime, error) {
	return &detector.Runtime{
		Name:    f.name,
		Version: f.version,
		Path:    f.path,
		Found:   f.found,
	}, nil
}

type fakeRuntimeClient struct {
	runtimes []*detector.Runtime
	err      error
	closed   bool
}

func (f *fakeRuntimeClient) ObserveRuntimes(ctx context.Context) ([]*detector.Runtime, error) {
	return f.runtimes, f.err
}

func (f *fakeRuntimeClient) Close() error {
	f.closed = true
	return nil
}
