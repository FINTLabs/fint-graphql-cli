package graphql

const SCHEMA_TEMPLATE = `# Built from tag {{ .GitTag }}

type {{ .Name }} {
{{- if .AttributesWithInheritance }}
	# Attributes
    {{- range $att := .AttributesWithInheritance }}
    {{ $att.Name }}: {{ graphqlType $att.Type | listFilt $att.List }}
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
