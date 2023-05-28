package entity

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Struct struct {
	name   string
	fields []*Field
}

func NewStruct(name string, fields []*Field) *Struct {
	return &Struct{
		name:   name,
		fields: fields,
	}
}

// 構造体のコンストラクタ関数のStatementを生成する
/*
	func <funcName>(
		<fieldValue[1]> <fieldType[1]>,
		<fieldValue[2]> <fieldType[2]>,
		...,
	) *<structName> {
		return &<structName>{
			<filedName[1]>: fieldValue[1],
			<fieldName[2]>: fieldValue[2],
			...,
		}
	}
*/
func (s *Struct) GenerateConstructorStatement() string {
	funcName := fmt.Sprintf("New%s", s.UpperName())

	stmt := jen.Func().Id(funcName).CustomFunc(
		jen.Options{
			Close:     ")",
			Multi:     true,
			Open:      "(",
			Separator: ",",
		},
		func(g *jen.Group) {
			for _, field := range s.fields {
				g.Id(field.LowerName()).Qual(field.typePkg, field.typeName)
			}
		}).Op("*").Id(s.name).Block(
		jen.Return(jen.Op("&").Id(s.name).CustomFunc(
			jen.Options{
				Close:     "}",
				Multi:     true,
				Open:      "{",
				Separator: ",",
			},
			func(g *jen.Group) {
				for _, field := range s.fields {
					g.Id(field.name).Op(":").Id(field.LowerName())
				}
			})),
	)

	return stmt.GoString()
}

// 構造体の各フィールドのGetterメソッドのStatementを生成する
/*
	func (x *<s.name>) <f.UpperName()>() <f.Type> {
		return x.<f.Name>
	}
*/
func (s *Struct) GenerateGettersStatement() []string {
	stmts := []string{}

	for _, f := range s.fields {
		if !f.NeedsGetter() {
			continue
		}

		stmt := jen.Func().Params(
			jen.Id("x").Op("*").Id(s.name),
		).Id(f.UpperName()).Params().Qual(f.typePkg, f.typeName).Block(
			jen.Return(jen.Id("x").Dot(f.name)),
		)

		stmts = append(stmts, stmt.GoString())
	}

	return stmts
}

// 構造体の名前の先頭を大文字に変換した文字列を返す。
func (s *Struct) UpperName() string {
	return cases.Title(language.Und, cases.NoLower).String(s.name)
}
