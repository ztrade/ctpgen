package main

import (
	"fmt"

	"github.com/go-clang/bootstrap/clang"
)

func WalkFile(file string, fn clang.CursorVisitor) (err error) {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()
	tu := idx.ParseTranslationUnit(file, []string{"-x", "c++"}, nil, 0)
	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		fmt.Println("PROBLEM:", d.Spelling())
	}
	cursor := tu.TranslationUnitCursor()
	cursor.Visit(fn)
	if len(diagnostics) > 0 {
		fmt.Println("NOTE: There were problems while analyzing the given file")
	}
	return
}
