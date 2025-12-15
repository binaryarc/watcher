package wsctl

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompletionCommandGeneratesBashScript(t *testing.T) {
	t.Helper()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "bash"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("completion command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "# bash completion for wsctl") {
		t.Fatalf("unexpected completion output: %s", output)
	}
}
