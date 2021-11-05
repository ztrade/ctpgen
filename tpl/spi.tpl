package {{.package}}

/*
#include "types_gen.h"
#include "{{.include}}"
*/
import "C"

type {{.className}} interface{
{{range $method := .methods}}
{{- $lenArg := len $method.Args -}}
{{$method.Name}}(
{{- range $i, $arg := $method.Args}}
     {{- $arg.Name}} {{$arg.Type}}
    {{- if isNotLast $i $lenArg -}}
        ,
    {{- end -}}
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
, {{- $arg.Name}} {{$arg.Type}}
  {{- end -}}){
    p := get{{$.className}}(ptr)
    if p != nil{
      p.{{$method.Name}}(
  {{- range $i, $arg := $method.Args}}
     {{- $arg.Name}}
     {{- if isNotLast $i $lenArg -}}
        ,
     {{- end -}}
  {{- end -}}
  )
  }
}
{{- end -}}