package graphql

const ENDPOINTS_TEMPLATE = `package no.fint.graphql.model;

import lombok.Getter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Getter
@Component
public class Endpoints {
{{range .}}
    @Value("${fint.endpoint.{{dots .}}:/{{.}}}")
    private String {{name .}};
{{end}}
}
`
