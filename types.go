package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-clang/bootstrap/clang"
)

var (
	//go:embed tpl/types.tpl
	TypesSource string

	//go:embed tpl/types_header.tpl
	TypesHeaderSource string
)

type Field struct {
	Name string
	Type string
}

type Struct struct {
	Name   string
	Fields []Field
}

type StructList struct {
	src   string
	Datas []*Struct
}

func goType(cType string) string {
	switch cType {
	case "int", "int64":
		return "int"
	case "double":
		return "float64"
	case "short":
		return "int16"
	case "char":
		return "byte"
	default:
		if strings.Contains(cType, "char [") {
			return "string"
		}
	}
	return cType
}

func (sl *StructList) Generate(pkg, dir string) (err error) {
	err = sl.GenerateCpp(pkg, dir)
	if err != nil {
		return
	}
	err = sl.GenerateGo(pkg, dir)
	return
}

func (sl *StructList) GenerateGo(pkg, dir string) (err error) {
	data := map[string]interface{}{
		"package": pkg,
		"structs": sl.Datas,
	}
	var fns template.FuncMap = map[string]interface{}{
		"goType":      goType,
		"goFieldName": goFieldName,
		"isCStr":      isCStr,
	}
	tmpl, err := template.New("types_go").Funcs(fns).Parse(TypesSource)
	if err != nil {
		return
	}
	file := filepath.Join(dir, "types_gen.go")
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		if err != nil {
			return
		}
		cmd := exec.Command("gofmt", "-w", file)
		buf, err1 := cmd.CombinedOutput()
		if err1 != nil {
			fmt.Println("gofmt error:", err1.Error())
		}
		fmt.Println(string(buf))
	}()
	err = tmpl.Execute(f, data)
	return
}

func (sl *StructList) GenerateCpp(pkg, dir string) (err error) {
	data := map[string]interface{}{
		"HeaderOnce": "_TYPES_GEN_H_",
		"src":        sl.src,
		"structs":    sl.Datas,
	}
	tmpl, err := template.New("types_header").Parse(TypesHeaderSource)
	if err != nil {
		return
	}
	file := filepath.Join(dir, "types_gen.h")
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	err = tmpl.Execute(f, data)
	return
}

func ParseStructData(structFile, dataFile string) (sl *StructList, err error) {
	sl = &StructList{src: structFile}
	sl.Datas, err = ParseStruct(structFile)
	if err != nil {
		return
	}
	dataTypes, err := ParseDataTypes(dataFile)
	if err != nil {
		return
	}
	var temp string
	var ok bool
	for _, v := range sl.Datas {
		for k, f := range v.Fields {
			temp, ok = dataTypes[f.Type]
			if ok {
				v.Fields[k].Type = temp
			}
		}
	}
	return
}

func ParseStruct(structFile string) (structs []*Struct, err error) {
	var st *Struct
	fn := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		switch cursor.Kind() {
		case clang.Cursor_StructDecl:
			if st != nil {
				structs = append(structs, st)
				st = nil
			}
			st = &Struct{}
			st.Name = cursor.Spelling()
			return clang.ChildVisit_Recurse
		case clang.Cursor_FieldDecl:
			if st == nil {
				return clang.ChildVisit_Continue
			}
			st.Fields = append(st.Fields, Field{Name: cursor.Spelling(), Type: cursor.Type().Spelling()})
		case clang.Cursor_EnumDecl, clang.Cursor_Namespace:
			return clang.ChildVisit_Continue
		}
		return clang.ChildVisit_Continue
	}
	err = WalkFile(structFile, fn)
	if err != nil {
		return
	}
	if st != nil && st.Name != "" {
		structs = append(structs, st)
	}

	return
}

func ParseDataTypes(dataFile string) (dataTypes map[string]string, err error) {
	dataTypes = make(map[string]string)
	fn := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		switch cursor.Kind() {
		case clang.Cursor_TypedefDecl:
			dataTypes[cursor.Spelling()] = cursor.TypedefDeclUnderlyingType().Spelling()
			return clang.ChildVisit_Continue
		}
		return clang.ChildVisit_Continue
	}
	err = WalkFile(dataFile, fn)
	if err != nil {
		return
	}
	return

}
