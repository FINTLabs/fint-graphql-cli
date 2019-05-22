package json

const SCHEMA_TEMPLATE = `{
    "$comment": "Version: {{.Tag}} generated at {{ timestamp }}",
    "$schema": "http://json-schema.org/schema#",
    "$id": "https://fintlabs.no/schema/{{ component .Package}}/{{ lowerCase .Name}}.json",
    "type": "object",
    "properties": {
{{- if .Attributes }}
    {{- range $att := .Attributes }}
        "{{ $att.Name }}": {
            "type": {{ $att | jsonType }}
        },
    {{- end }}
{{- end }}
{{- if .InheritedAttributes }}
    {{- range $att := .InheritedAttributes }}
        "{{ $att.Name }}": {
            "type": {{ $att | jsonTypeInh }}
        },
    {{- end }}
{{- end }}
        "_links": {
            "type": "object",
            "properties": {
{{- if .Relations -}}
	{{- range $i, $rel := . | jsonRelations }}
        {{- if $i }}, {{ end }}
                "{{ $rel.Name }}": {
                    "type": "array",
                    {{- if not $rel.Optional }}
                    "minItems": 1,
                    {{- end }}
                    {{- if not $rel.List }}
                    "maxItems": 1,
                    {{- end }}
                    "items": {
                        "type": "object",
                        "properties": {
                            "href": {
                                "type": "string",
                                "format": "uri"
                            }
                        },
                        "required": [ "href" ]
                    }
                }
	{{- end -}}
{{- end }}
            },
            "additionalProperties": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "href": {
                            "type": "string",
                            "format": "uri"
                        }
                    },
                    "required": [ "href" ]
                }
            },
            "required": [
{{- if .Relations -}}
    {{- range $i, $rel := . | requiredRelations }}
            {{- if $i}}, {{ end }}
                "{{ $rel }}"
	{{- end -}}
{{- end }}
            ]
        }
    },
    "required": [
        {{- range $i, $name := . | requiredAttributes -}}
        {{- if $i}}, {{end}}
        "{{$name}}"
        {{- end }}
    ]
}
`
