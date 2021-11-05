#ifndef {{.HeaderOnce}}
#define {{.HeaderOnce}}
#include <stdint.h>
#include "{{.src}}"
class {{.className}} : public {{.name}}
{
  {{.className}}(uint64_t ptr);
  virtual ~{{.className}}();
{{range $method := .methods}}
  virtual void {{$method.Name}}(
{{- $lenArg := len $method.Args -}}
{{- range $i, $arg := $method.Args}}
{{- $arg.Type}} {{$arg.Name}}
{{- if isNotLast $i $lenArg -}}
,
{{- end -}}
{{- end -}});
{{end}}

private:
  uint64_t ptr = 0;
};
#endif  // {{.HeaderOnce}}