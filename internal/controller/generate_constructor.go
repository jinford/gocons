package controller

import (
	"fmt"

	"github.com/jinford/gocons/internal/domain/service"
	"github.com/jinford/gocons/internal/repository"
)

func GenerateConstructor(parser service.Parser, generater service.CodeGenerater, printer repository.Printer) error {
	structs, err := parser.ParseSturcts()
	if err != nil {
		return fmt.Errorf("failed to pasre src structs: %w", err)
	}

	code, err := generater.GenerateCode(parser.SrcPkgName(), structs)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	if err := printer.Print(code); err != nil {
		return fmt.Errorf("failed to print generate code: %w", err)
	}

	return nil
}
