package json

const SCHEMA_TEMPLATE = `{
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
{{ if .InheritedAttributes }}
    {{- range $att := .InheritedAttributes }}
        "{{ $att.Name }}": {
            "type": {{ $att | jsonTypeInh }}
        },
    {{- end }}
{{- end }}
        "_links": {
            "type": "object",
            "properties": {
{{ if .Relations }}
	{{- range $i, $rel := .Relations }}
                "{{ $rel.Name }}": {
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
	{{ end -}}
{{- end }}
                "self": {
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
                }
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
    {{- range $i, $rel := .Relations }}
        {{- if not $rel.Optional }}
                "{{ $rel.Name }}",
        {{- end -}}
	{{- end -}}
{{- end }}
                "self"
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
