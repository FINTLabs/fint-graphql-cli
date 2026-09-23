package generate

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FINTLabs/fint-graphql-cli/common/types"
	"github.com/urfave/cli"
)

func mainClass(name, pkg string) *types.Class {
	return &types.Class{Name: name, Package: pkg, Stereotype: "hovedklasse", Identifiable: true,
		Identifiers: []types.Identifier{{Name: "systemId"}},
		Attributes:  []types.Attribute{{Name: "systemId", Type: "Identifikator"}}}
}

func requireContains(t *testing.T, output string, expected ...string) {
	t.Helper()
	for _, text := range expected {
		if !strings.Contains(output, text) {
			t.Errorf("missing %q in:\n%s", text, output)
		}
	}
}

func requireAbsent(t *testing.T, output string, unexpected ...string) {
	t.Helper()
	for _, text := range unexpected {
		if strings.Contains(output, text) {
			t.Errorf("unexpected %q in:\n%s", text, output)
		}
	}
}

func TestExplicitSpringQueryMappings(t *testing.T) {
	for _, name := range []string{"OtUngdom", "AvlagtProve", "Person"} {
		t.Run(name, func(t *testing.T) {
			c := mainClass(name, "no.novari.fint.model.utdanning.ot")
			c.Identifiers = append(c.Identifiers, types.Identifier{Name: "fodselsnummer", Optional: true})
			output := GetGraphQlQueryReolver(c)
			requireContains(t, output, "@Controller(\"model"+name+"QueryResolver\")",
				"@QueryMapping(name = \""+strings.ToLower(name)+"\")",
				"@Argument(\"systemId\") String systemId", "@Argument(\"fodselsnummer\") String fodselsnummer",
				"Resource> "+strings.ToLower(name)+"(", ".toFuture()",
				"import org.springframework.graphql.data.method.annotation.Argument;\nimport org.springframework.graphql.data.method.annotation.QueryMapping;",
				"import reactor.core.publisher.Mono;\n\nimport java.util.concurrent.CompletionStage;")
			requireAbsent(t, output, "kickstart", "coxautodev", "GraphQLQueryResolver", "@Component")
			if output != GetGraphQlQueryReolver(c) {
				t.Fatal("generation must be deterministic")
			}
		})
	}
}

func TestHiddenQueriesHaveNoMappingImports(t *testing.T) {
	c := mainClass("Funksjon", "no.novari.fint.model.administrasjon.kodeverk")
	requireAbsent(t, GetGraphQlQueryReolver(c), "QueryMapping", "Argument")
	c.Package = "no.novari.fint.model.utdanning.ot"
	c.Identifiers[0].Optional = true
	requireAbsent(t, GetGraphQlQueryReolver(c), "QueryMapping", "Argument")
	c.Identifiers = nil
	requireAbsent(t, GetGraphQlQueryReolver(c), "QueryMapping", "Argument", "StringUtils")
}

func TestRelationshipMappingsAndListSemantics(t *testing.T) {
	c := mainClass("OtUngdom", "no.novari.fint.model.utdanning.ot")
	c.Relations = []types.Association{
		{Name: "avlagtprove", Target: "AvlagtProve", TargetPackage: "no.novari.fint.model.utdanning.vurdering", Stereotype: "hovedklasse", List: true},
		{Name: "registrertAv", Target: "Skoleressurs", TargetPackage: "no.novari.fint.model.utdanning.elev", Stereotype: "hovedklasse"},
		{Name: "foreldre", Target: "OtUngdom", TargetPackage: c.Package, Stereotype: "hovedklasse", List: true},
	}
	output := GetGraphQlResolver(c)
	requireContains(t, output,
		"@SchemaMapping(typeName = \"OtUngdom\", field = \"avlagtprove\")",
		"@SchemaMapping(typeName = \"OtUngdom\", field = \"registrertAv\")",
		"getAvlagtprove(OtUngdomResource", "getRegistrertAv(OtUngdomResource",
		"flatMapSequential", "8, 1)", ".defaultIfEmpty(Optional.empty())",
		".onErrorResume(WebClientResponseException.class", "opt.orElse(null)",
		"CompletableFuture.completedFuture(List.of())", ".next()",
		"import java.util.List;\nimport java.util.Optional;\nimport java.util.concurrent.CompletableFuture;\nimport java.util.concurrent.CompletionStage;")
	requireAbsent(t, output, "GraphQLResolver", "coxautodev", "import no.fint.graphql.model.model.otungdom.OtUngdomService;")
	if strings.Count(output, "import no.novari.fint.model.resource.utdanning.ot.OtUngdomResource;") != 1 {
		t.Fatal("self-resource imports must be deduplicated")
	}
	c.Relations = c.Relations[1:2]
	requireAbsent(t, GetGraphQlResolver(c), "import java.util.List;", "Optional;", "CompletableFuture;", "Collectors;", "WebClientResponseException;")
}

func TestResourceWildcardThreshold(t *testing.T) {
	pkg := "no.novari.fint.model.resource.administrasjon.kodeverk."
	names := []string{pkg + "ArtResource", pkg + "AnleggResource", pkg + "AnsvarResource", pkg + "FunksjonResource", pkg + "RammeResource"}
	requireAbsent(t, formatImports(names[:4]), ".*;")
	output := formatImports(append(names, names[0]))
	if output != "import "+pkg+"*;" {
		t.Fatalf("unexpected imports: %s", output)
	}
}

func TestOptionalFeideIdentifiersAreNormalized(t *testing.T) {
	for _, name := range []string{"Elev", "Skoleressurs"} {
		t.Run(name, func(t *testing.T) {
			c := mainClass(name, "no.novari.fint.model.utdanning.elev")
			attribute := types.Attribute{Name: "feidenavn", Type: "Identifikator", Optional: true}
			c.Attributes = append(c.Attributes, attribute)
			requireContains(t, GetGraphQlService(c), "import org.springframework.util.StringUtils;",
				"resource.getFeidenavn() != null && !StringUtils.hasText(resource.getFeidenavn().getIdentifikatorverdi())", "resource.setFeidenavn(null)")
			c.Attributes = nil
			c.InheritedAttributes = []types.InheritedAttribute{{Attribute: attribute}}
			requireContains(t, GetGraphQlService(c), "resource.setFeidenavn(null)")
			c.InheritedAttributes[0].Optional = false
			requireAbsent(t, GetGraphQlService(c), "StringUtils", ".map(")
		})
	}
	c := mainClass("Person", "no.novari.fint.model.felles")
	requireAbsent(t, GetGraphQlService(c), "StringUtils", "getFeidenavn")
}

func TestValidScalarAndEmptyTypeGeneration(t *testing.T) {
	c := mainClass("Person", "no.novari.fint.model.felles")
	requireContains(t, GetGraphQlRootSchema([]*types.Class{c}), "scalar Date\nscalar Long", "person(systemId: String): Person")
	for _, name := range []string{"Grepreferanse", "Vigoreferanse"} {
		if GetGraphQlSchema(&types.Class{Name: name}) != "" {
			t.Errorf("empty type %s must not be emitted", name)
		}
	}
	c.Attributes = nil
	c.InheritedAttributes = []types.InheritedAttribute{{Attribute: types.Attribute{Name: "belop", Type: "long"}}}
	requireContains(t, GetGraphQlSchema(c), "belop: Long!")
}

func TestSchemaExclusionsAlsoExcludeQueryMappings(t *testing.T) {
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(previous)
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	cli.StringSliceFlag{Name: "exclude-schema"}.Apply(flags)
	if err := flags.Parse([]string{"--exclude-schema", "OtUngdom"}); err != nil {
		t.Fatal(err)
	}
	context := cli.NewContext(cli.NewApp(), flags, nil)
	classes := []*types.Class{mainClass("OtUngdom", "no.novari.fint.model.utdanning.ot"), mainClass("Person", "no.novari.fint.model.felles")}
	setupGraphQlSchemaDirStructure()
	generateGraphQlSchema(classes, context)
	generateGraphQlQueryResolver(classes, context)
	root, err := os.ReadFile(filepath.Join("graphql", "schema", "root.graphqls"))
	if err != nil {
		t.Fatal(err)
	}
	requireAbsent(t, string(root), "otungdom")
	controller, err := os.ReadFile(filepath.Join("graphql", "model", "model", "otungdom", "OtUngdomQueryResolver.java"))
	if err != nil {
		t.Fatal(err)
	}
	requireAbsent(t, string(controller), "QueryMapping", "Argument")
}
