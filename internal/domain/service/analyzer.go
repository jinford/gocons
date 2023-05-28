package service

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
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
	config := &types.Config{
		Importer: importer.Default(),
	}

	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}

	pkgName := a.astFile.Name.Name
	pkg, err := config.Check(pkgName, a.tokenFileset, []*ast.File{a.astFile}, info)
	if err != nil {
		return nil, fmt.Errorf("type-checks error: %w", err)
	}

	structs := []*entity.Struct{}
	ast.Inspect(a.astFile, func(node ast.Node) bool {
		t, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structName := t.Name.Name
		s, ok := pkg.Scope().Lookup(structName).Type().Underlying().(*types.Struct)
		if !ok {
			return true
		}

		fileds := a.parseFields(pkgName, s)
		structs = append(structs, entity.NewStruct(structName, fileds))

		return true
	})

	return structs, nil
}

func (a *astParser) parseFields(pkgName string, s *types.Struct) []*entity.Field {
	fs := make([]*entity.Field, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)

		typePkg, typeName, found := strings.Cut(field.Type().String(), ".")
		if !found {
			typePkg = ""
			typeName = field.Type().String()
		}

		if typePkg == a.SrcPkgName() {
			typePkg = ""
		}

		tagRawValue := reflect.StructTag(s.Tag(i)).Get(a.structTag)
		tagValues := strings.Split(tagRawValue, ",")

		fs[i] = entity.NewField(field.Name(), typePkg, typeName, field.Exported(), tagValues)
	}

	return fs
}
