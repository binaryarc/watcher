package grpcserver

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/proto"
)

type fakeDetector struct {
	name    string
	version string
	path    string
	found   bool
	calls   *int32
}

func (f *fakeDetector) Name() string {
	return f.name
}

func (f *fakeDetector) Detect() (*detector.Runtime, error) {
	if f.calls != nil {
		atomic.AddInt32(f.calls, 1)
	}
	return &detector.Runtime{
		Name:    f.name,
		Version: f.version,
		Path:    f.path,
		Found:   f.found,
	}, nil
}

func TestWatcherServerObserveRuntimes(t *testing.T) {
	t.Parallel()

	var detectorCalls int32
	detectorsProvider = func() []detector.Detector {
		return []detector.Detector{
			&fakeDetector{
				name:    "java",
				version: "17",
				path:    "/opt/java/bin/java",
				found:   true,
				calls:   &detectorCalls,
			},
			&fakeDetector{
				name:  "missing",
				found: false,
			},
		}
	}
	t.Cleanup(func() {
		detectorsProvider = detector.GetAllDetectors
	})

	server := NewWatcherServer()
	resp, err := server.ObserveRuntimes(context.Background(), &proto.ObserveRequest{})
	if err != nil {
		t.Fatalf("ObserveRuntimes() error = %v", err)
	}

	if len(resp.Runtimes) != 1 {
		t.Fatalf("expected 1 runtime, got %d", len(resp.Runtimes))
	}

	rt := resp.Runtimes[0]
	if rt.Name != "java" || rt.Version != "17" || rt.Path != "/opt/java/bin/java" {
		t.Fatalf("unexpected runtime: %+v", rt)
	}

	if resp.SystemInfo == nil || resp.SystemInfo.Hostname == "" {
		t.Fatalf("system info not populated: %+v", resp.SystemInfo)
	}

	if resp.Timestamp == 0 {
		t.Fatalf("timestamp not set")
	}

	if atomic.LoadInt32(&detectorCalls) == 0 {
		t.Fatalf("detector Detect() not called")
	}
}

func TestWatcherServerRuntimeFilter(t *testing.T) {
	t.Parallel()

	detectorsProvider = func() []detector.Detector {
		return []detector.Detector{
			&fakeDetector{name: "java", found: true},
			&fakeDetector{name: "python", found: true},
		}
	}
	t.Cleanup(func() {
		detectorsProvider = detector.GetAllDetectors
	})

	server := NewWatcherServer()
	resp, err := server.ObserveRuntimes(context.Background(), &proto.ObserveRequest{
		RuntimeFilter: []string{"python"},
	})
	if err != nil {
		t.Fatalf("ObserveRuntimes() error = %v", err)
	}

	if len(resp.Runtimes) != 1 || resp.Runtimes[0].Name != "python" {
		t.Fatalf("filter not applied, runtimes: %+v", resp.Runtimes)
	}
}

func TestFilterDetectors(t *testing.T) {
	t.Parallel()

	detectors := []detector.Detector{
		&fakeDetector{name: "java"},
		&fakeDetector{name: "python"},
	}

	filtered := filterDetectors(detectors, []string{"python"})
	if len(filtered) != 1 || filtered[0].Name() != "python" {
		t.Fatalf("filterDetectors() = %+v", filtered)
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	t.Parallel()

	now := time.Now().Unix()
	ts := getCurrentTimestamp()
	if ts < now || ts > now+5 {
		t.Fatalf("timestamp %d out of expected range", ts)
	}
}
