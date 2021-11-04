package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-clang/bootstrap/clang"
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
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()
	tu := idx.ParseTranslationUnit(file, []string{"-std=c++11"}, nil, 0)
	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		fmt.Println("PROBLEM:", d.Spelling())
	}
	cursor := tu.TranslationUnitCursor()
	var cls clang.Cursor
	spi = &SpiClass{Src: file}
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
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
	})
	if spi.Name == "" {
		err = fmt.Errorf("spi not found")
		return
	}
	if len(diagnostics) > 0 {
		fmt.Println("NOTE: There were problems while analyzing the given file")
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
