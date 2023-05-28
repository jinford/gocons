package repository

import (
	"fmt"
	"os"
)

type stdoutPrinter struct{}

var _ Printer = (*stdoutPrinter)(nil)

func NewStdoutPrinter() *stdoutPrinter {
	return &stdoutPrinter{}
}

func (p *stdoutPrinter) Print(code []byte) error {
	if _, err := fmt.Fprint(os.Stdout, string(code)); err != nil {
		return fmt.Errorf("failed to write generated code to stdout: %w", err)
	}
	return nil
}
