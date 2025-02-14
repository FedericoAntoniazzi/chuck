package printers

import (
	"fmt"
)

const requiredRowParamsCount = 3

type TextPrinter struct {
	rows []string
}

func (cp *TextPrinter) Initialize() {
	cp.rows = make([]string, 0)
}

func (cp *TextPrinter) AddRowParams(parts ...string) {
	rowParams := []string{"", "", ""}

	for i := 0; i < requiredRowParamsCount; i++ {
		if i < len(parts) {
			rowParams[i] = parts[i]
		}
	}

	row := fmt.Sprintf("Container %s (%s) can be upgraded to %s", rowParams[0], rowParams[1], rowParams[2])
	cp.rows = append(cp.rows, row)
}

func (cp *TextPrinter) Print() {
	for i := 0; i < len(cp.rows); i++ {
		fmt.Println(cp.rows[i])
	}
}
