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
	"graphqlType":    types.GetGraphQlType,
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
	"getPathFromPackage": func(p string, n string) string {
		//"no.fint.model.administrasjon.personal"
		s := strings.Replace(p, "no.fint.model.", "", -1)
		return fmt.Sprintf("%s/%s", strings.Replace(s, ".", "/", -1), strings.ToLower(n))
	},
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
	"graphqlRelation": func(t *types.Association) string {
		if t.Stereotype == "hovedklasse" {
			return t.Target
		}
		return "String"
	},
}

func GetGraphQlSchema(c *types.Class) string {
	return GetSchema(c, graphql.SCHEMA_TEMPLATE)
}

func GetGraphQlQueryReolver(c *types.Class) string {
	return getClass(c, graphql.QUERY_RESOLVER_TEMPLATE)
}

func GetGraphQlService(c *types.Class) string {
	return getClass(c, graphql.SERVICE_TEMPLATE)
}

func GetGraphQlResolver(c *types.Class) string {
	return getClass(c, graphql.RESOLVER_TEMPLATE)
}

func getClass(c *types.Class, t string) string {
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
