package graphql

const SERVICE_TEMPLATE = `// Built from tag {{ .GitTag }}

package no.fint.graphql.model.{{ lowerCase .Name}};

import no.fint.graphql.ResourceUrlBuilder;
import {{resourcePkg .Package}}.{{ .Name }}Resource;
import {{resourcePkg .Package}}.{{ .Name }}Resources;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

@Service
public class {{ .Name }}Service {

    @Autowired
    private WebClient webClient;

    public {{ .Name }}Resources get{{ .Name }}Resources(String sinceTimeStamp) {


        return webClient.get()
                .uri(ResourceUrlBuilder.urlWithQueryParams("{{ getPathFromPackage .Package .Name}}", sinceTimeStamp))
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
