package graphql

const SCHEMA_TEMPLATE = `
type {{ .Name }} {
{{- if .Attributes }}
	#  Attributes
    {{- range $att := .Attributes }}
    {{ $att.Name }}: {{ $att | graphqlType }}
    {{- end }}
{{- end }}
{{ if .InheritedAttributes }}
	# Inherited Attributes
    {{- range $att := .InheritedAttributes }}
    {{ $att.Name }}: {{ $att | graphqlTypeInh }} # {{ $att.Owner }}
    {{- end }}
{{- end }}
{{ if .Relations }}
	# Relations
	{{- range $i, $rel := .Relations }}
    {{ $rel.Name }}: {{ $rel | graphqlRelation }}
	{{- end }}
{{ end -}}
}
`
