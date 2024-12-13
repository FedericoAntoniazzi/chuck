package outputs

import (
	"fmt"
	"io"

	"github.com/FedericoAntoniazzi/chuck/internal/imagecheck"
)

type ConsoleOutput struct {
	Output io.Writer
}

func NewConsoleOutput(output io.Writer) ConsoleOutput {
	return ConsoleOutput{
		Output: output,
	}
}

func (co ConsoleOutput) Submit(results []imagecheck.UpdateResult) {
	for _, result := range results {
		if !result.HasUpdate {
			continue
		}

		fmt.Fprintf(co.Output, "Container %s can be updated to %s:%s\n", result.Container.Name, result.Container.Image, result.LatestVersion)
	}
}
