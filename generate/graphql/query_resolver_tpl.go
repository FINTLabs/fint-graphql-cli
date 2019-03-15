package graphql

const QUERY_RESOLVER_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component("{{component .Package}}{{.Name}}QueryResolver")
public class {{ .Name }}QueryResolver implements GraphQLQueryResolver {

    @Autowired
    private {{ .Name }}Service service;

    public {{ .Name }}Resource get{{ .Name }}(
{{- range $i, $ident := .Identifiers }}
            String {{ .Name }},
{{- end }}
            DataFetchingEnvironment dfe) {
{{- range $i, $ident := .Identifiers }}
        if (StringUtils.isNotEmpty({{ .Name }})) {
            return service.get{{ $.Name }}ResourceById("{{ lowerCase .Name }}", {{.Name}}, dfe);
        }
{{- end }}
        return null;
    }
}
`
