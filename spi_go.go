package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

var (
	//go:embed go.tmpl
	GoSource string
)

func goType(typ string) (goTyp string) {
	if typ == "" {
		return
	}
	if typ[len(typ)-1] == '*' {
		goTyp = "* C." + typ[0:len(typ)-1]
	} else {
		goTyp = typ
	}
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
			methods[k].Args[argi].Type = goType(arg.Type)
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
	tmpl, err := template.New("spigo").Funcs(fns).Parse(GoSource)
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
