package graphql

const ROOT_TEMPLATE = `# java.util.Date implementation
scalar Date

type Query {
{{- range . }}
	{{ lowerCase .Name }}({{range $i,$e := .Identifiers}}{{if $i}}, {{end}}{{.Name}}: String{{end}}): {{ .Name }}
{{- end }}
}
`
