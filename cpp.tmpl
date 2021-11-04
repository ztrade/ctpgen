#include "{{.Header}}"

extern "C"{
{{range $method := .methods}}
  void {{$.prefix}}{{$method.Name}}(uint64_t ptr
{{- range $i, $arg := $method.Args -}}
,{{- $arg.Type}} {{$arg.Name}}
{{- end -}});
{{end}}
}

{{.className}}::{{.className}}(uint64_t ptr):ptr(ptr){
}

{{.className}}::~{{.className}}(){
}

{{range $method := .methods}}
void {{$.className}}::{{$method.Name}}(
{{- $lenArg := len $method.Args -}}
{{- range $i, $arg := $method.Args}}
    {{- $arg.Type}} {{$arg.Name}}
    {{- if isNotLast $i $lenArg -}}
        ,
    {{- end -}}
{{- end -}})
{
  {{$.prefix}}{{$method.Name}}(ptr
{{- range $i, $arg := $method.Args -}}
,{{- $arg.Name}}
{{- end -}});

}
{{end}}