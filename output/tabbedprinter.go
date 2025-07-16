package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"go.uber.org/zap"
)

// TabbedPrinter abstracts a writer to print out tabstopped text into aligned text
type TabbedPrinter struct {
	writer *tabwriter.Writer
	logger *zap.SugaredLogger
}

// NewTabbedPrinter creates a default TabbedPrinter
func NewTabbedPrinter(log *zap.SugaredLogger) *TabbedPrinter {
	return &TabbedPrinter{
		writer: tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight),
		logger: log,
	}
}

// SetHeaders set titles for each column of the output
func (tp *TabbedPrinter) SetHeaders(headers ...string) {
	_, err := fmt.Fprintln(tp.writer, strings.Join(headers, "\t"))
	if err != nil {
		tp.logger.Errorf("error setting headers: %v", err)
	}
}

// AddRow add rows to the output
func (tp *TabbedPrinter) AddRow(columns ...any) {
	_, err := fmt.Fprintln(tp.writer, strings.Join(normalizeStrings(columns), "\t"))
	if err != nil {
		tp.logger.Errorf("error adding result row: %v", err)
	}
}

// Print flushes the writer buffer to default output
func (tp *TabbedPrinter) Print() {
	err := tp.writer.Flush()
	if err != nil {
		tp.logger.Errorf("error flushing buffer from tabbedprinter: %v", err)
	}
}

func normalizeStrings(items []any) []string {
	normalized := make([]string, len(items))

	for pos, item := range items {
		normalized[pos] = fmt.Sprint(item)
	}

	return normalized
}
