#ifndef {{.HeaderOnce}}
#define {{.HeaderOnce}}
#include "{{.src}}"

{{range $struct := .structs}}
typedef struct {{$struct.Name}} {{$struct.Name}};
{{- end }}

#endif  // {{.HeaderOnce}}
