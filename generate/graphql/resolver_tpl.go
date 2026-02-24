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

import no.novari.fint.model.resource.Link;
import {{resourcePkg .Package}}.{{ .Name}}Resource;
{{- if $ur -}}
{{- range $i, $rel := $ur }}
import {{ resourcePkg $rel.TargetPackage }}.{{ $rel.Target }}Resource;
{{- end -}}
{{- end }}

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClientResponseException;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletionStage;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

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
        {{ if $rel.List -}}
        var links = Optional.ofNullable({{ lowerCase $.Name}}.get{{ $rel.Name | upperCaseFirst }}()).orElseGet(List::of);
        if (links.isEmpty()) {
            return CompletableFuture.completedFuture(List.of());
        }
        return Flux.fromIterable(links)
                .map(Link::getHref)
                .flatMapSequential(href -> {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(href, dfe)
                        .map(Optional::of)
                        .onErrorResume(WebClientResponseException.NotFound.class,
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
