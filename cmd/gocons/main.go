package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinford/gocons/internal/controller"
	"github.com/jinford/gocons/internal/domain/service"
	"github.com/jinford/gocons/internal/repository"
	"github.com/urfave/cli/v2"
)

// These variables are set in build step
var (
	Version = "dev"
)

const (
	appName = "gocons"

	flagNameSrc    = "src"
	flagNameTag    = "tag"
	flagNameOutput = "output"
	flagNameValues = "values"
)

func main() {
	app := &cli.App{
		Name:    appName,
		Usage:   "generate constructor function & getter methods for structs",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagNameSrc,
				Usage:    "path of file that declares structs (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  flagNameTag,
				Usage: "name of target struct tag",
				Value: "cons",
			},
			&cli.StringFlag{
				Name:  flagNameOutput,
				Usage: "output: 'file', 'stdout'",
				Value: "file",
			},
			&cli.BoolFlag{
				Name:  flagNameValues,
				Usage: "generate constructor returning the value struct, instead of the pointer one",
			},
		},
		Action: execute,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(cCtx *cli.Context) error {
	src := cCtx.String(flagNameSrc)
	tag := cCtx.String(flagNameTag)
	output := cCtx.String(flagNameOutput)
	isValues := cCtx.Bool(flagNameValues)

	parser, err := service.NewParser(src, tag)
	if err != nil {
		return err
	}

	generater := service.NewCodeGenerator(appName, isValues)

	var printer repository.Printer
	switch output {
	case "file":
		ext := filepath.Ext(src)
		before, _ := strings.CutSuffix(src, ext)
		dest := fmt.Sprintf("%s.consgen%s", before, ext)
		printer = repository.NewFilePrinter(dest)
	case "stdout":
		printer = repository.NewStdoutPrinter()
	default:
		return fmt.Errorf("unkown 'output' option value inputted: '%s'", output)
	}

	if err := controller.GenerateConstructor(parser, generater, printer); err != nil {
		return err
	}

	return nil
}
