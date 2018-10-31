package graphql

const SERVICE_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ component .Package }}.{{ lowerCase .Name}};

import no.fint.graphql.model.Endpoints;
import no.fint.graphql.ResourceUrlBuilder;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import {{resourcePkg .Package}}.{{ .Name }}Resources;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

@Service("{{component .Package}}{{.Name}}Service")
public class {{ .Name }}Service {

    @Autowired
    private WebClient webClient;

    @Autowired
    private Endpoints endpoints;

    public {{ .Name }}Resources get{{ .Name }}Resources(String sinceTimeStamp) {


        return webClient.get()
                .uri(ResourceUrlBuilder.urlWithQueryParams(endpoints.{{ .Package | getPathFromPackage | getEndpoint }} + "/{{ lowerCase .Name }}", sinceTimeStamp))
                .retrieve()
                .bodyToMono({{ .Name }}Resources.class)
                .block();
    }

    public {{ .Name }}Resource get{{ .Name }}Resource(String url) {
        return webClient.get()
                .uri(url)
                .retrieve()
                .bodyToMono({{ .Name }}Resource.class)
                .block();
    }
}

`