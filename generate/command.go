package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/FINTLabs/fint-graphql-cli/common/config"
	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/FINTLabs/fint-graphql-cli/common/parser"
	"github.com/FINTLabs/fint-graphql-cli/common/types"
	"github.com/FINTLabs/fint-graphql-cli/common/utils"
	"github.com/urfave/cli"
)

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		fmt.Print("Getting latest from GitHub...")
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
		fmt.Printf(" %s\n", tag)
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")
	owner := c.GlobalString("owner")
	repo := c.GlobalString("repo")
	filename := c.GlobalString("filename")
	exclude := c.StringSlice("exclude")

	classesOrig, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)
	var classes []*types.Class
	if len(exclude) == 0 {
		classes = classesOrig
	} else {
		for _, class := range classesOrig {
			if !isStringInList(class.Name, exclude) {
				var rel []types.Association
				for _, r := range class.Relations {
					if !isStringInList(r.Target, exclude) {
						rel = append(rel, r)
					} else {
						fmt.Printf("Excluding %s.%s from %s.%s\n", r.TargetPackage, r.Target, class.Package, class.Name)
					}
				}
				class.Relations = rel
				var att []types.Attribute
				for _, a := range class.Attributes {
					if !isStringInList(a.Type, exclude) {
						att = append(att, a)
					} else {
						fmt.Printf("Excluding %s (%s) from %s.%s\n", a.Name, a.Type, class.Package, class.Name)
					}
				}
				class.Attributes = att
				classes = append(classes, class)
			} else {
				fmt.Printf("Excluding class %s.%s\n", class.Package, class.Name)
			}
		}
	}

	setupGraphQlSchemaDirStructure()
	generateGraphQlSchema(classes, c)
	generateGraphQlQueryResolver(classes)
	generateGraphQlService(classes)
	generateGraphQlResolver(classes)
	fmt.Println("Done!")
}

func writeFile(path string, filename string, content []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path+"/"+filename, content, 0777)
}

func getModelPath(pkg string, className string) string {
	return fmt.Sprintf("%s/model/%s/%s", config.GRAPHQL_BASE_PATH, types.GetComponentName(pkg), strings.ToLower(className))
}

func writeSchema(pkg string, schema string, content []byte) error {
	path := fmt.Sprintf("%s/schema/%s", config.GRAPHQL_BASE_PATH, types.GetComponentName(pkg))
	return writeFile(path, strings.ToLower(schema)+".graphqls", []byte(content))
}

func writeQueryResolver(pkg string, className string, content []byte) error {
	path := getModelPath(pkg, className)
	return writeFile(path, fmt.Sprintf("%sQueryResolver.java", className), []byte(content))
}

func writeService(pkg string, className string, content []byte) error {
	path := getModelPath(pkg, className)
	return writeFile(path, fmt.Sprintf("%sService.java", className), []byte(content))
}

func writeResolver(pkg string, className string, content []byte) error {
	path := getModelPath(pkg, className)
	return writeFile(path, fmt.Sprintf("%sResolver.java", className), []byte(content))
}

func generateGraphQlSchema(classes []*types.Class, cli *cli.Context) {

	fmt.Println("Generating GraphQL Schema")

	var roots []*types.Class

	for _, c := range classes {
		if !c.Abstract && includePackage(c.Package) && !excludeFromSchema(cli, c.Name) {
			fmt.Printf("  > Creating schema: %s.graphqls\n", c.Name)
			schema := GetGraphQlSchema(c)
			err := writeSchema(c.Package, c.Name, []byte(schema))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
			if c.Stereotype == "hovedklasse" && c.Identifiable && !strings.Contains(c.Package, "kodeverk") {
				include := false
				for _, i := range c.Identifiers {
					include = include || !i.Optional
				}
				if include {
					roots = append(roots, c)
				}
			}
		}
	}

	fmt.Println("  > Creating schema: root.graphqls")
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Name < roots[j].Name
	})

	schema := GetGraphQlRootSchema(roots)
	err := writeFile(config.GRAPHQL_BASE_PATH+"/schema", "root.graphqls", []byte(schema))
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}

func generateGraphQlQueryResolver(classes []*types.Class) {

	fmt.Println("Generating GraphQL Query Resolvers")

	for _, c := range classes {
		if !c.Abstract && c.Stereotype == "hovedklasse" && includePackage(c.Package) {
			fmt.Printf("  > Creating query resolver: %s.java\n", c.Name)
			class := GetGraphQlQueryReolver(c)
			err := writeQueryResolver(c.Package, c.Name, []byte(class))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}
}

func generateGraphQlService(classes []*types.Class) {

	fmt.Println("Generating GraphQL Services")

	var resources []string

	for _, c := range classes {
		if !c.Abstract && c.Stereotype == "hovedklasse" && includePackage(c.Package) {
			fmt.Printf("  > Creating service: %s.java\n", c.Name)
			class := GetGraphQlService(c)
			resources = append(resources, GetPackagePath(c.Package))
			err := writeService(c.Package, c.Name, []byte(class))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}

	fmt.Println("  > Creating Endpoints.java")
	class := GetEndpoints(utils.Distinct(resources))
	err := writeFile(config.GRAPHQL_BASE_PATH+"/model", "Endpoints.java", []byte(class))
	if err != nil {
		fmt.Printf("Unable to write file Endpoints.java: %s", err)
	}
}

func generateGraphQlResolver(classes []*types.Class) {

	fmt.Println("Generating GraphQL Resolvers")

	for _, c := range classes {
		if !c.Abstract && c.Resource && includePackage(c.Package) {
			fmt.Printf("  > Creating resolver: %s.java\n", c.Name)
			class := GetGraphQlResolver(c)
			err := writeResolver(c.Package, c.Name, []byte(class))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}
}

func setupGraphQlSchemaDirStructure() {
	os.RemoveAll("graphql")
	os.Mkdir(config.GRAPHQL_BASE_PATH, 0777)
}

func includePackage(p string) bool {
	return strings.Contains(p, "administrasjon") || strings.Contains(p, "utdanning") || strings.Contains(p, "felles")
}

func excludeFromSchema(c *cli.Context, p string) bool {
	excludeSchema := c.StringSlice("exclude-schema")

	if len(excludeSchema) == 0 {
		return false
	}

	for _, s := range excludeSchema {
		if strings.EqualFold(p, s) {
			return true
		}
	}
	return false
}

func isStringInList(s string, list []string) bool {
	for _, v := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}
