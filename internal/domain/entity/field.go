package entity

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Field struct {
	name      string
	typePkg   string
	typeName  string
	exported  bool
	tagValues []string
}

func NewField(name string, typePkg string, typeName string, exported bool, tagValues []string) *Field {
	return &Field{
		name:      name,
		typePkg:   typePkg,
		typeName:  typeName,
		exported:  exported,
		tagValues: tagValues,
	}
}

// Getterメソッドを生成する必要があるか判定する。
// フィールドがプライベートのとき、かつconsタグに"getter"が含まれているときの生成する必要がある。
func (f *Field) NeedsGetter() bool {
	if f.exported {
		return false
	}

	for _, v := range f.tagValues {
		if v == "getter" {
			return true
		}
	}

	return false
}

// フィールド名の先頭を小文字に変換した文字列を返す。
// ただし、フィールド名が全て大文字で構成されている場合は全て小文字に変換した文字列を返す。
func (f *Field) LowerName() string {
	if len(f.name) == 0 {
		return f.name
	}

	r := []rune(f.name)

	upperAll := make([]rune, len(r))
	for i, r := range f.name {
		upperAll[i] = unicode.ToUpper(r)
	}

	if string(upperAll) == f.name {
		return strings.ToLower(f.name)
	}

	r[0] = unicode.ToLower(r[0])

	return string(r)
}

// フィールド名の先頭を大文字に変換した文字列を返す。
func (f *Field) UpperName() string {
	return cases.Title(language.Und, cases.NoLower).String(f.name)
}
