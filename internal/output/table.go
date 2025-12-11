package output

import (
	"os"

	"github.com/binaryarc/watcher/internal/detector"
	"github.com/olekukonko/tablewriter"
)

// PrintRuntimesTable prints runtimes in table format
func PrintRuntimesTable(runtimes []*detector.Runtime) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Runtime", "Version", "Path"})

	for _, rt := range runtimes {
		if rt.Found {
			table.Append([]string{
				rt.Name,
				rt.Version,
				rt.Path,
			})
		}
	}

	table.Render()
}

// PrintRuntimeTable prints a single runtime in table format
func PrintRuntimeTable(runtime *detector.Runtime) {
	if !runtime.Found {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Property", "Value"})

	table.Append([]string{"Name", runtime.Name})
	table.Append([]string{"Version", runtime.Version})
	table.Append([]string{"Path", runtime.Path})

	table.Render()
}
