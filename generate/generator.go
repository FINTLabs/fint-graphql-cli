package generate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/FINTLabs/fint-graphql-cli/generate/graphql"

	"github.com/FINTLabs/fint-graphql-cli/common/types"
)

var funcMap = template.FuncMap{
	"normalizeFeidenavn": normalizeFeidenavn,
	"queryImports":       queryImports,
	"resolverImports":    resolverImports,
	//"add": func(i int, ii int) int { return i + ii },
	"sub": func(i int, ii int) int { return i - ii },
	"resourcePkg": func(s string) string {
		return strings.Replace(s, "model", "model.resource", -1)
	},
	"resource": func(resources []types.Attribute, s string) string {
		for _, a := range resources {
			if strings.HasSuffix(s, a.Type) {
				return strings.Replace(s, "model", "model.resource", -1) + "Resource"
			}
		}
		return s
	},
	"extends": func(isResource bool, extends string, s string) string {
		if isResource && strings.HasSuffix(s, extends) {
			return strings.Replace(s, "model", "model.resource", -1) + "Resource"
		}
		return s
	},
	"listFilt": func(list bool, s string) string {
		if list {
			return fmt.Sprintf("List<%s>", s)
		}
		return s
	},
	"javaType": types.GetJavaType,
	"csType": func(s string, opt bool) string {
		typ := types.GetCSType(s)
		if opt && types.IsValueType(typ) {
			return typ + "?"
		}
		return typ
	},
	"component":   types.GetComponentName,
	"graphqlType": types.GetGraphQlType,
	"graphqlTypeInh": func(a *types.InheritedAttribute) string {
		return types.GetGraphQlType(&a.Attribute)
	},
	"graphqlRelation": types.GetGraphQlRelationType,
	"relTargetType": func(a *types.Association) string {
		if a.List {
			return fmt.Sprintf("List<%sResource>", a.Target)
		}
		return a.Target + "Resource"
	},
	"lowerCase":      func(s string) string { return strings.ToLower(s) },
	"upperCase":      func(s string) string { return strings.ToUpper(s) },
	"upperCaseFirst": func(s string) string { return strings.Title(s) },
	"getter":         func(s string) string { return "get" + strings.Title(s) + "()" },
	"baseType":       func(s string) string { return strings.Replace(s, "Resource", "", -1) },
	"assignResource": func(typ string, att string) string {
		if strings.HasPrefix(typ, "List<") {
			inner := strings.TrimSuffix(strings.TrimPrefix(typ, "List<"), ">")
			return fmt.Sprintf("%s.stream().map(%s::create).collect(Collectors.toList())", att, inner)
		}
		return fmt.Sprintf("%s.create(%s)", typ, att)
	},
	"listAdder": func(typ string) string {
		if strings.HasPrefix(typ, "List<") {
			return "All"
		}
		return ""
	},
	"getPathFromPackage": GetPackagePath,
	"uniqueRelationTargets": func(input []types.Association) []types.Association {
		u := make([]types.Association, 0, len(input))
		m := make(map[string]bool)

		for _, val := range input {
			if val.Stereotype == "hovedklasse" {
				if _, ok := m[val.Target]; !ok {
					m[val.Target] = true
					u = append(u, val)
				}
			}
		}

		return u
	},
	"getEndpoint": func(r string) string { return "get" + strings.Title(GetEndpointName(r)) + "()" },
	"endpointForClass": func(c *types.Class) string {
		return "get" + strings.Title(GetEndpointName(GetEndpointPathForClass(c))) + "()"
	},
}

func GetPackagePath(p string) string {
	return strings.Join(strings.Split(p, ".")[4:], "/")
}

func GetEndpointName(p string) string {
	var r string
	for i, s := range strings.Split(p, "/") {
		if i == 0 {
			r += s
		} else {
			r += strings.Title(s)
		}
	}
	return r
}

func GetEndpointPathForClass(c *types.Class) string {
	return GetPackagePath(c.Package)
}

func GetGraphQlSchema(c *types.Class) string {
	if len(c.Attributes) == 0 && len(c.InheritedAttributes) == 0 && len(c.Relations) == 0 {
		return ""
	}
	return GetSchema(c, graphql.SCHEMA_TEMPLATE)
}

func GetGraphQlQueryReolver(c *types.Class) string {
	return getQueryResolver(c, isRootQuery(c))
}

type queryResolverModel struct {
	*types.Class
	PublicQuery bool
}

func getQueryResolver(c *types.Class, publicQuery bool) string {
	return getClass(&queryResolverModel{Class: c, PublicQuery: publicQuery}, graphql.QUERY_RESOLVER_TEMPLATE)
}

// Share eligibility between SDL roots and annotated controllers.
func isRootQuery(c *types.Class) bool {
	if c.Abstract || c.Stereotype != "hovedklasse" || !c.Identifiable || strings.Contains(c.Package, "kodeverk") || !includePackage(c.Package) {
		return false
	}
	for _, identifier := range c.Identifiers {
		if !identifier.Optional {
			return true
		}
	}
	return false
}

func GetGraphQlService(c *types.Class) string {
	return getClass(c, graphql.SERVICE_TEMPLATE)
}

// An empty optional Feide identifier must be absent, not an invalid Identifikator.
func normalizeFeidenavn(c *types.Class) bool {
	attributes := append([]types.Attribute{}, c.Attributes...)
	for _, inherited := range c.InheritedAttributes {
		attributes = append(attributes, inherited.Attribute)
	}
	for _, attribute := range attributes {
		if attribute.Name == "feidenavn" && attribute.Type == "Identifikator" && attribute.Optional && !attribute.List {
			return true
		}
	}
	return false
}

func GetGraphQlResolver(c *types.Class) string {
	return getClass(c, graphql.RESOLVER_TEMPLATE)
}

func GetGraphQlRootSchema(classes []*types.Class) string {
	tpl := template.New("schema").Funcs(funcMap)

	parse, err := tpl.Parse(graphql.ROOT_TEMPLATE)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, classes)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func GetEndpoints(r []string) string {
	var funcs = template.FuncMap{
		"dots": func(s string) string { return strings.Replace(s, "/", ".", -1) },
		"name": GetEndpointName,
	}
	tpl := template.New("class").Funcs(funcs)

	parse, err := tpl.Parse(graphql.ENDPOINTS_TEMPLATE)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, r)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func getClass(c interface{}, t string) string {
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, c)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func GetSchema(c *types.Class, t string) string {
	tpl := template.New("schema").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, c)
	if err != nil {
		panic(err)
	}
	return b.String()
}
