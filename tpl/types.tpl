// This file was automatically generated by ctpgen
package {{.package}}

/*
#include "gen_types.h"
*/
import "C"


type THOST_TE_RESUME_TYPE int
var(
	THOST_TERT_RESTART THOST_TE_RESUME_TYPE = 0
	THOST_TERT_RESUME THOST_TE_RESUME_TYPE = 1
	THOST_TERT_QUICK THOST_TE_RESUME_TYPE = 2
)

{{range $struct := .structs}}
type {{$struct.Name}} struct{
  {{- range $field := $struct.Fields }}
  {{goFieldName $field.Name}} {{goType $field.Type}}
  {{- end -}}
}

func New{{$struct.Name}}(p *C.{{$struct.Name}}) *{{$struct.Name}}{
  ret := new({{$struct.Name}})
  {{range $field := $struct.Fields -}}
    {{if isCStr $field.Type -}}
  ret.{{$field.Name}} = c2goStr(&p.{{$field.Name}}[0], C.sizeof_{{$struct.Name}})
    {{else -}}
  ret.{{$field.Name}} = {{goType $field.Type}}(p.{{$field.Name}})
    {{end -}}
  {{end -}}
  return ret
}

func {{$struct.Name}}CValue(s *{{$struct.Name}}) *C.{{$struct.Name}} {
  ptr := (*C.{{$struct.Name}})(C.malloc(C.sizeof_{{$struct.Name}}))
  {{range $field := $struct.Fields -}}
    {{if isCStr $field.Type -}}
    go2cStr(s.{{$field.Name}}, &ptr.{{$field.Name}}[0], C.sizeof_{{$struct.Name}})
    {{else -}}
    ptr.{{$field.Name}} = C.{{$field.Type}}(s.{{$field.Name}})
    {{end -}}
  {{end -}}
  return ptr
}

{{end}}