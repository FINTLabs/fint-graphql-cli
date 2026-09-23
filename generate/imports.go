package generate

import (
	"sort"
	"strings"

	"github.com/FINTLabs/fint-graphql-cli/common/types"
)

func resourcePackage(pkg string) string {
	return strings.Replace(pkg, ".model.", ".model.resource.", 1)
}

// Derive imports from the template model instead of rewriting generated Java.
func queryImports(c *queryResolverModel) string {
	names := []string{
		"graphql.schema.DataFetchingEnvironment",
		"lombok.extern.slf4j.Slf4j",
		resourcePackage(c.Package) + "." + c.Name + "Resource",
		"org.springframework.beans.factory.annotation.Autowired",
		"org.springframework.stereotype.Controller",
		"reactor.core.publisher.Mono",
		"java.util.concurrent.CompletionStage",
	}
	if len(c.Identifiers) > 0 {
		names = append(names, "org.apache.commons.lang3.StringUtils")
	}
	if c.PublicQuery {
		names = append(names, "org.springframework.graphql.data.method.annotation.QueryMapping")
		if len(c.Identifiers) > 0 {
			names = append(names, "org.springframework.graphql.data.method.annotation.Argument")
		}
	}
	return formatImports(names)
}

func resolverImports(c *types.Class) string {
	names := []string{"org.springframework.stereotype.Controller"}
	hasRelations, hasLists := false, false
	for _, rel := range c.Relations {
		if rel.Stereotype != "hovedklasse" {
			continue
		}
		hasRelations = true
		hasLists = hasLists || rel.List
		if rel.Target != c.Name || rel.TargetPackage != c.Package {
			names = append(names, "no.fint.graphql.model."+types.GetComponentName(rel.TargetPackage)+"."+strings.ToLower(rel.Target)+"."+rel.Target+"Service")
		}
		names = append(names, resourcePackage(rel.TargetPackage)+"."+rel.Target+"Resource")
	}
	if hasRelations {
		names = append(names,
			"graphql.schema.DataFetchingEnvironment",
			"no.novari.fint.model.resource.Link",
			resourcePackage(c.Package)+"."+c.Name+"Resource",
			"org.springframework.beans.factory.annotation.Autowired",
			"org.springframework.graphql.data.method.annotation.SchemaMapping",
			"reactor.core.publisher.Flux", "reactor.core.publisher.Mono",
			"java.util.concurrent.CompletionStage",
		)
	}
	if hasLists {
		names = append(names,
			"org.springframework.web.reactive.function.client.WebClientResponseException",
			"java.util.List", "java.util.Optional", "java.util.concurrent.CompletableFuture",
			"java.util.stream.Collectors",
		)
	}
	return formatImports(names)
}

func formatImports(names []string) string {
	unique := map[string]bool{}
	resources := map[string][]string{}
	for _, name := range names {
		if unique[name] {
			continue
		}
		unique[name] = true
		if strings.HasPrefix(name, "no.novari.fint.model.resource.") && strings.HasSuffix(name, "Resource") {
			pkg := name[:strings.LastIndex(name, ".")]
			resources[pkg] = append(resources[pkg], name)
		}
	}
	// Match the maintained resource import threshold without service wildcards.
	for pkg, classes := range resources {
		if len(classes) >= 5 {
			for _, name := range classes {
				delete(unique, name)
			}
			unique[pkg+".*"] = true
		}
	}
	ordinary, java := []string{}, []string{}
	for name := range unique {
		line := "import " + name + ";"
		if strings.HasPrefix(name, "java.") {
			java = append(java, line)
		} else {
			ordinary = append(ordinary, line)
		}
	}
	sort.Strings(ordinary)
	sort.Strings(java)
	result := strings.Join(ordinary, "\n")
	if len(java) > 0 {
		result += "\n\n" + strings.Join(java, "\n")
	}
	return result
}
