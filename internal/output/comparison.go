package output

import (
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// ComparisonData represents the comparison result
type ComparisonData struct {
	Hosts    []string            `json:"hosts" yaml:"hosts"`
	Runtimes []RuntimeComparison `json:"runtimes" yaml:"runtimes"`
}

// RuntimeComparison represents a single runtime across all servers
type RuntimeComparison struct {
	Name     string   `json:"name" yaml:"name"`
	Versions []string `json:"versions" yaml:"versions"`
	Status   string   `json:"status" yaml:"status"`
}

// PrintComparisonTable prints runtime comparison in table format
func PrintComparisonTable(comparison *ComparisonData) {
	table := tablewriter.NewWriter(os.Stdout)

	// Build header: Runtime | Host1 | Host2 | ... | Status
	header := []string{"Runtime"}
	header = append(header, comparison.Hosts...)
	header = append(header, "Status")
	table.Header(header)

	// Add rows
	for _, rt := range comparison.Runtimes {
		row := []string{rt.Name}
		row = append(row, rt.Versions...)
		row = append(row, formatStatus(rt.Status))
		table.Append(row)
	}

	table.Render()
}

// formatStatus adds emoji to status
func formatStatus(status string) string {
	switch status {
	case "SAME":
		return "✅ SAME"
	case "DIFF":
		return "⚠️  DIFF"
	case "PARTIAL":
		return "⚠️  PARTIAL"
	case "MISSING":
		return "❌ MISSING"
	case "ERROR":
		return "❌ ERROR"
	default:
		return status
	}
}

// PrintComparisonJSON prints comparison in JSON format
func PrintComparisonJSON(comparison *ComparisonData) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(comparison)
}

// PrintComparisonYAML prints comparison in YAML format
func PrintComparisonYAML(comparison *ComparisonData) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(comparison)
}
