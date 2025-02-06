package outputs

import "fmt"

type ConsoleOutput struct{}

func NewConsoleOutput() ConsoleOutput {
	return ConsoleOutput{}
}

func (o *ConsoleOutput) Send(message string) {
	fmt.Println(message)
}
