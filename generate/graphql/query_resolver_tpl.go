package graphql

const QUERY_RESOLVER_TEMPLATE = `
package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

{{ queryImports . }}

@Controller("{{component .Package}}{{.Name}}QueryResolver")
@Slf4j
public class {{ .Name }}QueryResolver {

    @Autowired
    private {{ .Name }}Service service;

{{ if .PublicQuery }}    @QueryMapping(name = "{{ lowerCase .Name }}")
{{ end }}    public CompletionStage<{{ .Name }}Resource> {{ lowerCase .Name }}(
{{- range $i, $ident := .Identifiers }}
            {{ if $.PublicQuery }}@Argument("{{ .Name }}") {{ end }}String {{ .Name }},
{{- end }}
            DataFetchingEnvironment dfe) {
		log.info("New Query for {{ .Name }}");
{{- range $i, $ident := .Identifiers }}
        if (StringUtils.isNotEmpty({{ .Name }})) {
            return service.get{{ $.Name }}ResourceById("{{ lowerCase .Name }}", {{.Name}}, dfe).toFuture();
        }
{{- end }}
        return Mono.<{{ .Name }}Resource>empty().toFuture();
    }
}
`
