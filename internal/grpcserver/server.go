package grpcserver

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/binaryarc/watcher/proto"
)

type WatcherServer struct {
	proto.UnimplementedWatcherServiceServer
}

func NewWatcherServer() *WatcherServer {
	return &WatcherServer{}
}

func (s *WatcherServer) ObserveRuntimes(ctx context.Context, req *proto.ObserveRequest) (*proto.ObserveResponse, error) {
	// 1. 모든 detector 가져오기
	detectors := detector.GetAllDetectors()

	// 2. 필터가 있으면 적용
	if len(req.RuntimeFilter) > 0 {
		detectors = filterDetectors(detectors, req.RuntimeFilter)
	}

	// 3. 런타임 감지
	var protoRuntimes []*proto.Runtime
	for _, det := range detectors {
		runtime, err := det.Detect()
		if err != nil {
			continue
		}

		if runtime.Found {
			protoRuntimes = append(protoRuntimes, &proto.Runtime{
				Name:    runtime.Name,
				Version: runtime.Version,
				Path:    runtime.Path,
				Found:   runtime.Found,
			})
		}
	}

	// 4. 시스템 정보 수집
	hostname, _ := os.Hostname()
	systemInfo := &proto.SystemInfo{
		Hostname: hostname,
		Os:       getOS(),
		Kernel:   getKernel(),
	}

	// 5. 응답 생성
	response := &proto.ObserveResponse{
		Runtimes:   protoRuntimes,
		SystemInfo: systemInfo,
		Timestamp:  getCurrentTimestamp(),
	}

	return response, nil
}

func filterDetectors(detectors []detector.Detector, filters []string) []detector.Detector {
	filterMap := make(map[string]bool)
	for _, f := range filters {
		filterMap[f] = true
	}

	var filtered []detector.Detector
	for _, det := range detectors {
		if filterMap[det.Name()] {
			filtered = append(filtered, det)
		}
	}
	return filtered
}

func getOS() string {
	return runtime.GOOS
}

func getKernel() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
