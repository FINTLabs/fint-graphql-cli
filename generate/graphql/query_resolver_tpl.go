package graphql

const QUERY_RESOLVER_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import {{resourcePkg .Package}}.{{ .Name }}Resources;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class {{ .Name }}QueryResolver implements GraphQLQueryResolver {

    @Autowired
    private {{ .Name }}Service service;

    public List<{{ .Name }}Resource> get{{ .Name }}(String sinceTimeStamp) {
        {{ .Name }}Resources resources = service.get{{ .Name }}Resources(sinceTimeStamp);
        return resources.getContent();
    }
}
`
