package detector

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type JavaDetector struct{}

// Name은 detector 이름을 반환
func (d *JavaDetector) Name() string {
	return "java"
}

// Detect는 Java 버전을 감지
func (d *JavaDetector) Detect() (*Runtime, error) {
	runtime := &Runtime{
		Name:  "java",
		Found: false,
	}

	// 1. java 명령어가 있는지 확인
	javaPath, err := exec.LookPath("java")
	if err != nil {
		// java가 없으면 Found=false로 반환
		return runtime, nil
	}

	runtime.Path = javaPath
	runtime.Found = true

	// 2. java -version 실행
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput() // stderr로 출력되므로 CombinedOutput 사용
	if err != nil {
		return runtime, fmt.Errorf("failed to execute java -version: %w", err)
	}

	// 3. 버전 파싱
	version := parseJavaVersion(string(output))
	runtime.Version = version

	return runtime, nil
}

// parseJavaVersion은 java -version 출력에서 버전을 추출
func parseJavaVersion(output string) string {
	// java -version 출력 예시:
	// openjdk version "11.0.19" 2023-04-18
	// java version "1.8.0_372"
	// openjdk version "17.0.8" 2023-07-18 LTS

	// 정규식으로 버전 추출
	re := regexp.MustCompile(`version "(.+?)"`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		version := matches[1]
		// "1.8.0_372" 같은 경우 "8.0" 형태로 변환
		if strings.HasPrefix(version, "1.") {
			parts := strings.Split(version, ".")
			if len(parts) >= 2 {
				return parts[1] + ".x" // "1.8.0" -> "8.x"
			}
		}
		return version
	}

	return "unknown"
}
