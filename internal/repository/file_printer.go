package repository

import (
	"fmt"
	"os"
)

type filePrinter struct {
	dest string
}

var _ Printer = (*filePrinter)(nil)

func NewFilePrinter(dest string) *filePrinter {
	return &filePrinter{
		dest: dest,
	}
}

func (p *filePrinter) Print(code []byte) error {
	if err := os.WriteFile(p.dest, code, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", p.dest, err)
	}
	return nil
}
