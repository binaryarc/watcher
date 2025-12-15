package detector

import "testing"

func TestParseJavaVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{
			name:   "modern version",
			output: `openjdk version "17.0.8" 2023-07-18 LTS`,
			want:   "17.0.8",
		},
		{
			name:   "legacy 1.x version",
			output: `java version "1.8.0_372"`,
			want:   "8.x",
		},
		{
			name:   "unknown format",
			output: `some invalid output`,
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseJavaVersion(tt.output); got != tt.want {
				t.Fatalf("parseJavaVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParsePythonVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"valid", "Python 3.10.12", "3.10.12"},
		{"invalid", "something else", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parsePythonVersion(tt.output); got != tt.want {
				t.Fatalf("parsePythonVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseNodeVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"with prefix", "v20.18.0", "20.18.0"},
		{"without prefix", "18.12.1", "18.12.1"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseNodeVersion(tt.output); got != tt.want {
				t.Fatalf("parseNodeVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseGoVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"standard", "go version go1.21.5 linux/amd64", "1.21.5"},
		{"unknown", "invalid output", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseGoVersion(tt.output); got != tt.want {
				t.Fatalf("parseGoVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseDockerVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"valid", "Docker version 24.0.5, build ced0996", "24.0.5"},
		{"invalid", "Docker build unknown", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseDockerVersion(tt.output); got != tt.want {
				t.Fatalf("parseDockerVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseMySQLVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"mysql", "mysql  Ver 8.0.34 for Linux on x86_64", "8.0.34"},
		{"mariadb", "mysql  Ver 10.11.4-MariaDB", "10.11.4"},
		{"unknown", "mysql Distrib foo", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseMySQLVersion(tt.output); got != tt.want {
				t.Fatalf("parseMySQLVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseRedisVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"valid", "Redis server v=7.0.12 sha=000 malloc=jemalloc", "7.0.12"},
		{"invalid", "Redis info missing version", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseRedisVersion(tt.output); got != tt.want {
				t.Fatalf("parseRedisVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseNginxVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"valid", "nginx version: nginx/1.24.0", "1.24.0"},
		{"invalid", "nginx something", "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseNginxVersion(tt.output); got != tt.want {
				t.Fatalf("parseNginxVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestGetAllDetectors(t *testing.T) {
	t.Parallel()
	detectors := GetAllDetectors()
	wantNames := []string{"java", "python", "node", "go", "docker", "mysql", "redis", "nginx"}

	if len(detectors) != len(wantNames) {
		t.Fatalf("GetAllDetectors() = %d detectors, want %d", len(detectors), len(wantNames))
	}

	for i, det := range detectors {
		if det.Name() != wantNames[i] {
			t.Fatalf("detector[%d].Name() = %q, want %q", i, det.Name(), wantNames[i])
		}
	}
}
