package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func cTypeInGo(typ string) (goTyp string) {
	if typ == "" {
		return
	}
	typ = strings.Replace(typ, "bool", "C.int8_t", -1)
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
	if strings.Contains(typ, "char [") || strings.Contains(typ, "char[") {
		return true
	}
	return false
}

func goMethod(prefix string) string {
	return strings.ToUpper(prefix[0:1]) + prefix[1:]
}

// cMethod class method to c style
func cMethod(prefix, method string) string {
	var cMethod bytes.Buffer
	cMethod.WriteString(prefix)
	cMethod.WriteByte('_')
	buf := []byte(method)
	for k, v := range buf {
		if v < 'a' {
			if k != 0 {
				cMethod.WriteByte('_')
			}
			v += 'a' - 'A'
		}
		cMethod.WriteByte(v)

	}
	return cMethod.String()
}

func isSameType(typ1, typ2 string) bool {
	typ1 = strings.Replace(typ1, "*", "", -1)
	typ2 = strings.Replace(typ2, "*", "", -1)
	typ1 = strings.Replace(typ1, " ", "", -1)
	typ2 = strings.Replace(typ2, " ", "", -1)
	return typ1 == typ2
}

// cType
func cType(typ, name string) string {
	typ = strings.Replace(typ, "const", "", -1)
	typ = strings.Replace(typ, "bool", "int8_t", -1)
	// typ = strings.Replace(typ, "THOST_TE_RESUME_TYPE", "int", -1)
	typ = strings.Trim(typ, " ")
	if typ == "bool" {
		typ = "int8_t"
	}
	if typ == "char *[]" {
		return fmt.Sprintf("char* %s[]", name)
	}
	return fmt.Sprintf("%s %s", typ, name)
}

// goType ctype -> goType
func goType(cType string) string {
	cType = strings.Replace(cType, "const", "", -1)
	cType = strings.Trim(cType, " ")
	switch cType {
	case "void":
		return ""
	case "int", "int64":
		return "int"
	case "double":
		return "float64"
	case "short":
		return "int16"
	case "char":
		return "byte"
	default:
		if cType == "char *[]" {
			return "[]string"
		}
		if strings.Contains(cType, "char [") || strings.Contains(cType, "char[") || strings.Contains(cType, "char *") {
			return "string"
		}
		if strings.HasSuffix(cType, "Spi *") {
			cType = cType[0 : len(cType)-2]
			return cType
		}
		if cType[len(cType)-1] == '*' {
			cType = "* " + cType[0:len(cType)-2]
		}
	}
	return cType
}

// goToC go value -> c value
func goToC(typ string) (goTyp string) {
	if typ == "" {
		return
	}
	typ = strings.Replace(typ, "const", "", -1)
	typ = strings.Trim(typ, " ")
	switch typ {
	case "char *":
		return "go2cStrPtr"
	case "bool":
		return "go2cBool"
	case "char *[]":
		return "go2cStrArray"
	// case "THOST_TE_RESUME_TYPE":
	// 	return "C.enum_THOST_TE_RESUME_TYPE"
	default:
		if strings.HasPrefix(typ, "CThost") {
			return typ[0:len(typ)-2] + "CValue"
		}
		return fmt.Sprintf("C.%s", typ)
	}

}

// cToGo c value -> go value
func cToGo(typ, prefix, name string) string {
	name = prefix + name
	if typ == "" {
		return name
	}
	goTyp := strings.Replace(typ, "const", "", -1)
	goTyp = strings.Trim(goTyp, " ")
	switch goTyp {
	case "char *":
		goTyp = "cPtr2GoStr"
	case "bool":
		goTyp = "c2goBool"
	case "void":
		return ""
	case "int", "int64":
		goTyp = "int"
	case "double":
		goTyp = "goFloat64"
	case "short":
		goTyp = "int16"
	case "char":
		goTyp = "byte"
	default:
		if strings.Contains(goTyp, "char [") || strings.Contains(goTyp, "char[") {
			lenStr := strLen(goTyp)
			return fmt.Sprintf("c2goStr(&%s[0], %d)", name, lenStr)
		}
		if strings.HasPrefix(typ, "CThost") {
			goTyp = fmt.Sprintf("New%s", typ[0:len(typ)-2])
		} else if strings.HasSuffix(goTyp, "Spi *") {
			goTyp = goTyp[0 : len(goTyp)-2]
		} else if goTyp[len(goTyp)-1] == '*' {
			goTyp = "* " + goTyp[0:len(goTyp)-2]
		}
	}
	return fmt.Sprintf("%s(%s)", goTyp, name)
}

func isStrArray(args []Argument) bool {
	if len(args) < 2 {
		return false
	}
	if args[0].Type == "char *[]" && args[1].Type == "int" {
		return true
	}
	return false
}

func isNeedFree(typ string) bool {
	if typ == "char *" || typ == "char *[]" {
		return true
	}
	if strings.Contains(typ, "Spi *") {
		return false
	}
	if typ[len(typ)-1] == '*' {
		return true
	}
	return false
}

func freeMethod(typ, prefix, arg string) string {
	switch typ {
	case "char *":
		return fmt.Sprintf("freeCStr(%s%s)", prefix, arg)
	case "char *[]":
		return fmt.Sprintf("freeCStrArray(%s%s)", prefix, arg)
	}
	if typ[len(typ)-1] == '*' {
		return fmt.Sprintf("C.free(unsafe.Pointer(%s%s))", prefix, arg)
	}
	return ""
}

func strLen(typ string) int {
	typ = strings.Replace(typ, "char [", "", 1)
	typ = strings.Replace(typ, "char[", "", 1)
	typ = strings.Replace(typ, "]", "", 1)
	typ = strings.Trim(typ, " ")
	n, err := strconv.Atoi(typ)
	if err != nil {
		fmt.Println("atoi failed:", typ)
		return 1
	}
	return n
}
