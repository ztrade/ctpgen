// This file was automatically generated by ctpgen
#ifndef {{.HeaderOnce}}
#define {{.HeaderOnce}}
#include <stdint.h>
#include "{{.include}}"
#ifdef  __cplusplus
extern "C"{
#endif

typedef void* {{.prefix}}Api;
typedef void* {{.prefix}}Spi;

{{.prefix}}Spi {{.prefix}}_new_spi(uint64_t value);
void {{.prefix}}_spi_free({{.prefix}}Spi p);
{{range $method := .static_methods}}
{{if isSameType $method.Ret $.className}}
{{$.prefix}}Api {{cMethod $method.Name}}(
  {{- range $i,$arg := $method.Args -}}
  {{- if ne $i 0 -}}
  ,
  {{- end -}}
  {{cType $arg.Type $arg.Name}}
  {{- end -}}
);
{{- else -}}
{{$method.Ret}} {{cMethod $method.Name}}(
  {{- range $i,$arg := $method.Args -}}
  {{- if ne $i 0 -}}
  ,
  {{- end -}}
  {{cType $arg.Type $arg.Name}}
  {{- end -}}
);
{{- end}}
{{- end}}
{{range $method := .methods}}
{{if eq $method.Name "RegisterSpi" -}}
{{$method.Ret}} {{cMethod $method.Name}}({{$.prefix}}Api a, {{$.prefix}}Spi s);
{{- else -}}
{{$method.Ret}} {{cMethod $method.Name}}({{$.prefix}}Api a
  {{- range $arg := $method.Args -}}
  , {{cType $arg.Type $arg.Name}}
  {{- end -}}
);
{{- end}}
{{- end}}

#ifdef  __cplusplus
};
#endif
#endif  // {{.HeaderOnce}}