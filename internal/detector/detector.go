package detector

// Runtime 정보를 담는 구조체
type Runtime struct {
	Name    string // 프로그램 이름 (예: "java")
	Version string // 버전 (예: "11.0.19")
	Path    string // 실행 파일 경로 (예: "/usr/bin/java")
	Found   bool   // 발견 여부
}

// Detector 인터페이스
type Detector interface {
	Detect() (*Runtime, error)
	Name() string
}