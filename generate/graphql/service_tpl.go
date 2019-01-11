package graphql

const SERVICE_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import graphql.schema.DataFetchingEnvironment;
import no.fint.graphql.ResourceUrlBuilder;
import no.fint.graphql.WebClientRequest;
import no.fint.graphql.model.Endpoints;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import {{resourcePkg .Package}}.{{ .Name }}Resources;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service("{{component .Package}}{{.Name}}Service")
public class {{ .Name }}Service {

    @Autowired
    private WebClientRequest webClientRequest;

    @Autowired
    private Endpoints endpoints;

    public {{ .Name }}Resources get{{ .Name }}Resources(String sinceTimeStamp, DataFetchingEnvironment dfe) {
        return webClientRequest.get(
                ResourceUrlBuilder.urlWithQueryParams(
                    endpoints.{{ .Package | getPathFromPackage | getEndpoint }} + "/{{ lowerCase .Name }}",
                    sinceTimeStamp),
                {{ .Name }}Resources.class,
                dfe);
    }

    public {{ .Name }}Resource get{{ .Name }}Resource(String url, DataFetchingEnvironment dfe) {
        return webClientRequest.get(url, {{ .Name }}Resource.class, dfe);
    }
}

`
