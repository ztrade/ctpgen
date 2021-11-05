package main

import "strings"

func goTypeStyle(typ string) (goTyp string) {
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

func goFieldName(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func isCStr(typ string) bool {
	if strings.Contains(typ, "char [") {
		return true
	}
	return false
}
