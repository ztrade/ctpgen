package main

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
)

var (
	//go:embed header.tmpl
	CppHeader string
	//go:embed cpp.tmpl
	CppSource string
)

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
	tmpl, err := template.New("spi").Funcs(fns).Parse(CppHeader)
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
	tmpl, err := template.New("spi").Funcs(fns).Parse(CppSource)
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
