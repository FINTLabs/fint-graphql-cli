package graphql

const RESOLVER_TEMPLATE = `
package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

{{ resolverImports . }}

@Controller("{{component .Package}}{{.Name}}Resolver")
public class {{ .Name }}Resolver {
{{ $ur := uniqueRelationTargets .Relations -}}
{{ if $ur -}}
{{- range $i, $rel := $ur }}
    @Autowired
    private {{ $rel.Target }}Service {{ lowerCase $rel.Target}}Service;
{{ end }}
{{ end -}}

{{ range $i, $rel := .Relations -}}
{{- if eq $rel.Stereotype "hovedklasse" }}
    @SchemaMapping(typeName = "{{ $.Name }}", field = "{{ $rel.Name }}")
    public CompletionStage<{{ $rel | relTargetType }}> get{{ $rel.Name | upperCaseFirst }}({{ $.Name }}Resource {{ lowerCase $.Name}}, DataFetchingEnvironment dfe) {
        {{ if $rel.List -}}
        var links = Optional.ofNullable({{ lowerCase $.Name}}.get{{ $rel.Name | upperCaseFirst }}()).orElseGet(List::of);
        if (links.isEmpty()) {
            return CompletableFuture.completedFuture(List.of());
        }
        return Flux.fromIterable(links)
                .map(Link::getHref)
                .flatMapSequential(href -> {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(href, dfe)
                        .map(Optional::of)
                        .defaultIfEmpty(Optional.empty())
                        .onErrorResume(WebClientResponseException.class,
                                ex -> Mono.just(Optional.empty())),
                        8, 1)
                .collectList()
                .map(list -> list.stream()
                        .map(opt -> opt.orElse(null))
                        .collect(Collectors.toList()))
                .toFuture();
        {{- else -}}
        return Flux.fromStream({{ lowerCase $.Name}}.get{{ $rel.Name | upperCaseFirst }}()
                .stream()
                .map(Link::getHref)
                .map(l -> {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(l, dfe)))
                .flatMap(Mono::flux)
                .next()
                .toFuture();
        {{- end }}
    }
{{ end -}}
{{- end }}
}

`
