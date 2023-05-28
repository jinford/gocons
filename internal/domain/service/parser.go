package service

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinford/gocons/internal/domain/entity"
)

type Parser interface {
	SrcPkgName() string
	ParseSturcts() ([]*entity.Struct, error)
}

type astParser struct {
	src       string
	structTag string

	tokenFileset *token.FileSet
	astFile      *ast.File
}

func NewParser(src string, structTag string) (Parser, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse src file: %w", err)
	}

	return &astParser{
		src:          src,
		structTag:    structTag,
		tokenFileset: fset,
		astFile:      f,
	}, nil
}

func (a *astParser) SrcPkgName() string {
	return a.astFile.Name.Name
}

func (a *astParser) ParseSturcts() ([]*entity.Struct, error) {
	pkgName := a.astFile.Name.Name
	structs := []*entity.Struct{}
	ast.Inspect(a.astFile, func(node ast.Node) bool {
		t, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		s, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structName := t.Name.Name

		fileds := a.parseFields(pkgName, s)
		structs = append(structs, entity.NewStruct(structName, fileds))

		return true
	})

	return structs, nil
}

func (a *astParser) parseFields(pkgName string, s *ast.StructType) []*entity.Field {
	fs := make([]*entity.Field, len(s.Fields.List))
	for i, field := range s.Fields.List {
		fieldType := types.ExprString(field.Type)

		var fieldName string
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		} else {
			// embedding field
			_, typeName, found := strings.Cut(fieldType, ".")
			if !found {
				typeName = fieldType
			}

			// trim "*" as it could be a pointer
			fieldName = strings.TrimPrefix(typeName, "*")
		}

		tagValues := []string{}
		if field.Tag != nil {
			tag, err := strconv.Unquote(field.Tag.Value)
			if err != nil {
				panic(err)
			}

			tagRawValue := reflect.StructTag(tag).Get(a.structTag)
			tagValues = strings.Split(tagRawValue, ",")
		}

		fs[i] = entity.NewField(fieldName, fieldType, ast.IsExported(fieldName), tagValues)
	}

	return fs
}
