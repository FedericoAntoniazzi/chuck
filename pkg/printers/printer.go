package printers

import "errors"

type Printer interface {
	Initialize()
	AddRowParams(...string)
	Print()
}

func GetPrinter(kind string) (Printer, error) {
	if kind == "text" {
		return &TextPrinter{}, nil
	}
	if kind == "tabbed" {
		return &TabbedPrinter{}, nil
	}

	return nil, errors.New("Unknown printer type")
}
