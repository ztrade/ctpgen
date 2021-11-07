// This file was automatically generated by ctpgen
package {{.package}}

/*
#include <stdint.h>
#include "{{.include}}"
*/
import "C"

type {{.className}} interface{
{{range $method := .methods}}
{{- $lenArg := len $method.Args -}}
{{$method.Name}}(
{{- range $i, $arg := $method.Args}}
    {{- if ne $i 0 -}}
    ,
    {{- end -}}
    {{- $arg.Name}} {{goType $arg.Type}}
{{- end -}}
)
{{ end }}
}


func get{{.className}}(ptr uint64) {{.className}}{
     value := getGoPtr(ptr)
     if value == nil {
        return nil
     }
     v := value.({{.className}})
     return v
}


{{range $method := .methods}}
{{- $lenArg := len $method.Args}}
//export {{$.prefix}}{{$method.Name}}
func {{$.prefix}}{{$method.Name}}(ptr uint64
  {{- range $i, $arg := $method.Args -}}
, {{- $arg.Name}} {{cTypeInGo $arg.Type}}
  {{- end -}}){
    p := get{{$.className}}(ptr)
    if p != nil{
     {{- range $i, $arg := $method.Args}}
     go{{- $arg.Name}} := {{cToGo $arg.Type}}({{$arg.Name}})
     {{- end}}
      p.{{$method.Name}}(
  {{- range $i, $arg := $method.Args}}
     {{- if ne $i 0 -}}
     ,
     {{- end -}}
     go{{- $arg.Name}}
  {{- end -}}
  )
  }
}
{{- end -}}