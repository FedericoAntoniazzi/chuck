package outputs

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// TabbedPrinter abstracts a writer to print out tabstopped text into aligned text
type TabbedConsoleOutput struct {
	writer *tabwriter.Writer
}

// NewTabbedPrinter creates a default TabbedPrinter
func NewTabbedConsoleOutput() *TabbedConsoleOutput {
	return &TabbedConsoleOutput{
		writer: tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight),
	}
}

// SetHeaders set titles for each column of the output
func (tp *TabbedConsoleOutput) SetHeaders(headers ...string) {
	fmt.Fprintln(tp.writer, strings.Join(headers, "\t"))
}

// AddRow add rows to the output
func (tp *TabbedConsoleOutput) AddRow(columns ...any) {
	fmt.Fprintln(tp.writer, strings.Join(normalizeStrings(columns), "\t"))
}

// Print flushes the writer buffer to default output
func (tp *TabbedConsoleOutput) Print() {
	tp.writer.Flush()
}

func normalizeStrings(items []any) []string {
	normalized := make([]string, len(items))

	for pos, item := range items {
		normalized[pos] = fmt.Sprint(item)
	}

	return normalized
}
