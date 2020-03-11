package graphql

const RESOLVER_TEMPLATE = `
package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLResolver;
import graphql.schema.DataFetchingEnvironment;

{{ $ur := uniqueRelationTargets .Relations -}}
{{- if $ur -}}
{{- range $i, $rel := $ur -}}
import no.fint.graphql.model.{{ component $rel.TargetPackage }}.{{ lowerCase $rel.Target}}.{{ upperCaseFirst $rel.Target }}Service;
{{ end -}}
{{- end }}

import no.fint.model.resource.Link;
import {{resourcePkg .Package}}.{{ .Name}}Resource;
{{- if $ur -}}
{{- range $i, $rel := $ur }}
import {{ resourcePkg $rel.TargetPackage }}.{{ $rel.Target }}Resource;
{{- end -}}
{{- end }}

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.concurrent.CompletionStage;

@Component("{{component .Package}}{{.Name}}Resolver")
public class {{ .Name }}Resolver implements GraphQLResolver<{{ .Name }}Resource> {
{{ if $ur -}}
{{- range $i, $rel := $ur }}
    @Autowired
    private {{ $rel.Target }}Service {{ lowerCase $rel.Target}}Service;
{{ end }}
{{ end -}}

{{ range $i, $rel := .Relations -}}
{{- if eq $rel.Stereotype "hovedklasse" }}
    public CompletionStage<{{ $rel | relTargetType }}> get{{ $rel.Name | upperCaseFirst }}({{ $.Name }}Resource {{ lowerCase $.Name}}, DataFetchingEnvironment dfe) {
        return Flux.fromStream({{ lowerCase $.Name}}.get{{ $rel.Name | upperCaseFirst }}()
                .stream()
                .map(Link::getHref)
                .map(l -> {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(l, dfe)))
                .flatMap(Mono::flux)
                {{ if $rel.List -}}
                    .collectList()
                {{- else -}}
                    .next()
                {{- end }}
                .toFuture();
    }
{{ end -}}
{{- end }}
}

`
