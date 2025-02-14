package printers

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type TabbedPrinter struct {
	writer *tabwriter.Writer
}

func (tp *TabbedPrinter) Initialize() {
	tp.writer = tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

	headers := []string{"CONTAINER", "IMAGE", "VERSION UPDATE"}
	fmt.Fprintln(tp.writer, strings.Join(headers, "\t"))
}

func (tp *TabbedPrinter) AddRowParams(rows ...string) {
	fmt.Fprintln(tp.writer, strings.Join(rows, "\t"))
}

func (tp *TabbedPrinter) Print() {
	tp.writer.Flush()
}
