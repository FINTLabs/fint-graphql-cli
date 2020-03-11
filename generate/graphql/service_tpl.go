package graphql

const SERVICE_TEMPLATE = `
package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import graphql.schema.DataFetchingEnvironment;
import no.fint.graphql.WebClientRequest;
import no.fint.graphql.model.Endpoints;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

@Service("{{component .Package}}{{.Name}}Service")
public class {{ .Name }}Service {

    @Autowired
    private WebClientRequest webClientRequest;

    @Autowired
    private Endpoints endpoints;

    public Mono<{{ .Name }}Resource> get{{ .Name }}ResourceById(String id, String value, DataFetchingEnvironment dfe) {
        return get{{ .Name }}Resource(
            endpoints.{{ .Package | getPathFromPackage | getEndpoint }} 
                + "/{{ lowerCase .Name }}/" 
                + id 
                + "/" 
                + value,
            dfe);
    }

    public Mono<{{ .Name }}Resource> get{{ .Name }}Resource(String url, DataFetchingEnvironment dfe) {
        return webClientRequest.get(url, {{ .Name }}Resource.class, dfe);
    }
}

`
