package graphql

const QUERY_RESOLVER_TEMPLATE = `
package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import lombok.extern.slf4j.Slf4j;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.util.concurrent.CompletionStage;

@Component("{{component .Package}}{{.Name}}QueryResolver")
@Slf4j
public class {{ .Name }}QueryResolver implements GraphQLQueryResolver {

    @Autowired
    private {{ .Name }}Service service;

    public CompletionStage<{{ .Name }}Resource> get{{ .Name }}(
{{- range $i, $ident := .Identifiers }}
            String {{ .Name }},
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
