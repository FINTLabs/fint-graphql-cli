package graphql

const SCHEMA_TEMPLATE = `# Built from tag {{ .GitTag }}

type {{ .Name }} {
{{- if .AttributesWithInheritance }}
	# Attributes
    {{- range $att := .AttributesWithInheritance }}
    {{ $att.Name }}: {{ graphqlType $att }}
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
