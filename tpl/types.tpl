package {{.package}}

/*
#include "types_gen.h"
*/
import "C"

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

func (s *{{$struct.Name}}) CValue() *C.{{$struct.Name}} {
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