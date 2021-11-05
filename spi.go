package main

import (
	"bytes"
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
	//go:embed tpl/spi_header.tpl
	SpiCppHeader string
	//go:embed tpl/spi_cpp.tpl
	SpiCppSource string
	//go:embed tpl/spi.tpl
	SpiGoSource string
)

type Argument struct {
	Name string
	Type string
}

type ClassMethod struct {
	Name string
	Args []Argument
}

type SpiClass struct {
	Name    string
	Src     string
	Methods []ClassMethod
}

func ParseSpi(file string) (spi *SpiClass, err error) {
	var cls clang.Cursor
	spi = &SpiClass{Src: file}
	fn := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			fmt.Printf("cursor: <none>\n")
			return clang.ChildVisit_Continue
		}
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl:
			if strings.HasSuffix(cursor.Spelling(), "Spi") {
				cls = cursor
				spi.Name = cls.Spelling()
				return clang.ChildVisit_Recurse
			}
			return clang.ChildVisit_Continue
		case clang.Cursor_CXXMethod:
			if parent != cls {
				return clang.ChildVisit_Continue
			}
			method := ClassMethod{Name: cursor.Spelling()}
			for i := int32(0); i != cursor.NumArguments(); i++ {
				arg := cursor.Argument(uint32(i))
				method.Args = append(method.Args, Argument{Name: arg.Spelling(), Type: arg.Type().Spelling()})
			}
			spi.Methods = append(spi.Methods, method)
			return clang.ChildVisit_Continue
		case clang.Cursor_EnumDecl, clang.Cursor_StructDecl, clang.Cursor_Namespace:
			return clang.ChildVisit_Continue
		}

		return clang.ChildVisit_Continue
	}
	err = WalkFile(file, fn)
	if err != nil {
		return
	}
	if spi.Name == "" {
		err = fmt.Errorf("spi not found")
		return
	}
	return
}

func (s *SpiClass) Generate(pkg, prefix, dir string) (err error) {
	headerFile := prefix + "_spi.h"
	header := filepath.Join(dir, headerFile)
	err = s.GenerateCppHeader(header)
	if err != nil {
		err = fmt.Errorf("generate cpp header failed:%w", err)
		return
	}
	cppFile := prefix + "_spi.cpp"
	cpp := filepath.Join(dir, cppFile)
	err = s.GenerateCppSource(prefix, headerFile, cpp)
	if err != nil {
		err = fmt.Errorf("generate cpp source failed:%w", err)
		return
	}
	goFile := prefix + "_spi.go"
	goF := filepath.Join(dir, goFile)
	err = s.GenerateGo(prefix, pkg, goF)
	if err != nil {
		err = fmt.Errorf("generate go source failed:%w", err)
	}
	return
}

func isNotLast(i, total int) bool {
	return i != (total - 1)
}
func toUpperMacro(str string) string {
	buf := bytes.ToUpper([]byte(str))
	for k, v := range buf {
		if v == '.' {
			buf[k] = '_'
		}
	}
	return string(buf)
}
func (s *SpiClass) GenerateCppHeader(file string) (err error) {
	fName := filepath.Base(file)

	clsName := s.Name + "Impl"
	headerData := map[string]interface{}{
		"HeaderOnce": toUpperMacro(fName),
		"src":        s.Src,
		"className":  clsName,
		"name":       s.Name,
		"methods":    s.Methods,
	}
	var fns template.FuncMap = map[string]interface{}{
		"isNotLast": isNotLast,
	}
	tmpl, err := template.New("spi").Funcs(fns).Parse(SpiCppHeader)
	if err != nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	err = tmpl.Execute(f, headerData)
	return
}

func (s *SpiClass) GenerateCppSource(prefix, header, file string) (err error) {
	clsName := s.Name + "Impl"
	headerData := map[string]interface{}{
		"Header":    header,
		"src":       s.Src,
		"className": clsName,
		"name":      s.Name,
		"methods":   s.Methods,
		"prefix":    prefix,
	}
	var fns template.FuncMap = map[string]interface{}{
		"isNotLast": isNotLast,
	}
	tmpl, err := template.New("spi").Funcs(fns).Parse(SpiCppSource)
	if err != nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	err = tmpl.Execute(f, headerData)
	return
}

func (s *SpiClass) GenerateGo(prefix, pkg, file string) (err error) {
	clsName := s.Name
	methods := make([]ClassMethod, len(s.Methods))
	for k, v := range s.Methods {
		methods[k].Name = v.Name
		methods[k].Args = make([]Argument, len(v.Args))
		for argi, arg := range v.Args {
			methods[k].Args[argi].Name = arg.Name
			methods[k].Args[argi].Type = goTypeStyle(arg.Type)
		}
	}
	headerData := map[string]interface{}{
		"package":   pkg,
		"include":   "types_gen.h",
		"className": clsName,
		"name":      s.Name,
		"methods":   methods,
		"prefix":    prefix,
	}
	var fns template.FuncMap = map[string]interface{}{
		"isNotLast": isNotLast,
	}
	tmpl, err := template.New("spigo").Funcs(fns).Parse(SpiGoSource)
	if err != nil {
		return
	}

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
	err = tmpl.Execute(f, headerData)
	return
}
