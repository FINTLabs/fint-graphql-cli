package graphql

const RESOLVER_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLResolver;
import no.fint.graphql.Links;

{{ $ur := uniqueRelationTargets .Relations}}

{{ if $ur }}
{{- range $i, $rel := $ur }}
import no.fint.graphql.model.{{ component $rel.TargetPackage }}.{{ lowerCase $rel.Target}}.{{ upperCaseFirst $rel.Target }}Service;
{{- end }}
{{ end }}

import {{resourcePkg .Package}}.{{ .Name}}Resource;

{{ if $ur }}
{{- range $i, $rel := $ur }}
import {{ resourcePkg $rel.TargetPackage }}.{{ $rel.Target }}Resource;
{{- end }}
{{ end }}

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

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
    public {{ $rel.Target}}Resource get{{ $rel.Name | upperCaseFirst }}({{ $.Name }}Resource {{ lowerCase $.Name}}) {
        return {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(Links.get({{ lowerCase $.Name}}.get{{ $rel.Name | upperCaseFirst }}()));
    }
{{ end -}}
{{- end }}
}

`
