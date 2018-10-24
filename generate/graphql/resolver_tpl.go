package graphql

const RESOLVER_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLResolver;
import no.fint.graphql.Links;

{{ $ur := uniqueRelationTargets .Relations}}

{{ if $ur }}
{{- range $i, $rel := $ur }}
import no.fint.graphql.model.{{ lowerCase $rel.Target}}.{{ upperCaseFirst $rel.Target }}Service;
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

@Component
public class {{ .Name }}Resolver implements GraphQLResolver<{{ .Name }}Resource> {

	{{ if $ur }}
	{{- range $i, $rel := $ur }}
	@Autowired
	private {{ $rel.Target }}Service {{ lowerCase $rel.Target}}Service;
	{{ end }}
	{{ end }}


	{{ if $ur }}
	{{ $mainName := .Name }}
	{{- range $i, $rel := $ur }}
	public {{ $rel.Target}}Resource get{{ $rel.Target}}({{ $mainName }}Resource {{ lowerCase $mainName}}) {
        return {{ lowerCase $rel.Target}}Service.get{{ $rel.Target}}Resource(Links.get({{ lowerCase $mainName}}, "{{ $rel.Name }}"));
    }
	{{ end }}
	{{ end }}
}

`
