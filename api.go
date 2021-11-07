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
	//go:embed tpl/api_header.tpl
	ApiCppHeader string
	//go:embed tpl/api_cpp.tpl
	ApiCppSource string
	//go:embed tpl/api.tpl
	ApiGoSource string
)

type ApiClass struct {
	Name          string
	Src           string
	StaticMethods []ClassMethod
	Methods       []ClassMethod
	spiName       string
	SpiImplFile   string
}

func ParseApi(file string) (api *ApiClass, err error) {
	var cls clang.Cursor
	api = &ApiClass{Src: file}
	fn := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			fmt.Printf("cursor: <none>\n")
			return clang.ChildVisit_Continue
		}
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl:
			if strings.HasSuffix(cursor.Spelling(), "Api") {
				cls = cursor
				api.Name = cls.Spelling()
				return clang.ChildVisit_Recurse
			}
			if strings.HasSuffix(cursor.Spelling(), "Spi") {
				api.spiName = cursor.Spelling()
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
			method.Ret = cursor.ResultType().Spelling()
			if cursor.CXXMethod_IsStatic() {
				api.StaticMethods = append(api.StaticMethods, method)
			} else {
				api.Methods = append(api.Methods, method)
			}
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
	if api.Name == "" {
		err = fmt.Errorf("spi not found")
		return
	}
	return
}

func (s *ApiClass) Generate(pkg, prefix, dir string) (err error) {
	headerFile := "gen_" + prefix + "_api.h"
	header := filepath.Join(dir, headerFile)
	err = s.GenerateCppHeader(prefix, header)
	if err != nil {
		err = fmt.Errorf("generate cpp header failed:%w", err)
		return
	}
	cppFile := "gen_" + prefix + "_api.cpp"
	cpp := filepath.Join(dir, cppFile)
	err = s.GenerateCppSource(prefix, headerFile, cpp)
	if err != nil {
		err = fmt.Errorf("generate cpp source failed:%w", err)
		return
	}
	goFile := "gen_" + prefix + "_api.go"
	goF := filepath.Join(dir, goFile)
	err = s.GenerateGo(prefix, pkg, goF)
	if err != nil {
		err = fmt.Errorf("generate go source failed:%w", err)
	}
	return
}

func (s *ApiClass) GenerateCppHeader(prefix, file string) (err error) {
	fName := filepath.Base(file)
	headerData := map[string]interface{}{
		"HeaderOnce":     toUpperMacro(fName),
		"include":        "gen_types.h",
		"prefix":         prefix,
		"methods":        s.Methods,
		"static_methods": s.StaticMethods,
		"className":      s.Name,
	}
	var fns template.FuncMap = map[string]interface{}{
		"cMethod":    func(method string) string { return cMethod(prefix, method) },
		"isSameType": isSameType,
		"cType":      cType,
	}
	tmpl, err := template.New("api").Funcs(fns).Parse(ApiCppHeader)
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

func (s *ApiClass) GenerateCppSource(prefix, header, file string) (err error) {
	clsName := s.Name
	headerData := map[string]interface{}{
		"Header":         header,
		"include":        s.Src,
		"className":      clsName,
		"methods":        s.Methods,
		"static_methods": s.StaticMethods,
		"prefix":         prefix,
		"spiName":        s.spiName,
		"spiFile":        fmt.Sprintf("gen_%s_spi.h", prefix),
	}
	var fns template.FuncMap = map[string]interface{}{
		"cMethod":    func(method string) string { return cMethod(prefix, method) },
		"isSameType": isSameType,
		"cType":      cType,
	}
	tmpl, err := template.New("api_cpp").Funcs(fns).Parse(ApiCppSource)
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

func (s *ApiClass) GenerateGo(prefix, pkg, file string) (err error) {
	clsName := s.Name
	methods := make([]ClassMethod, len(s.Methods))
	staticMethods := []ClassMethod{}
	for _, v := range s.StaticMethods {
		staticMethods = append(staticMethods, v)
	}
	for k, v := range s.Methods {
		methods[k].Name = v.Name
		methods[k].Ret = v.Ret
		methods[k].Args = make([]Argument, len(v.Args))
		for argi, arg := range v.Args {
			methods[k].Args[argi].Name = arg.Name
			methods[k].Args[argi].Type = arg.Type
		}
	}
	headerData := map[string]interface{}{
		"package":        pkg,
		"include":        fmt.Sprintf("gen_%s_api.h", prefix),
		"name":           clsName,
		"methods":        methods,
		"prefix":         prefix,
		"static_methods": staticMethods,
		"spiName":        s.spiName,
	}
	var fns template.FuncMap = map[string]interface{}{
		"cMethod":    func(method string) string { return cMethod(prefix, method) },
		"goType":     goType,
		"goToC":      goToC,
		"cToGo":      cToGo,
		"isSameType": isSameType,
		"isStrArray": isStrArray,
		"goMethod":   goMethod,
		"isNeedFree": isNeedFree,
		"freeMethod": freeMethod,
	}
	tmpl, err := template.New("api_go").Funcs(fns).Parse(ApiGoSource)
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
