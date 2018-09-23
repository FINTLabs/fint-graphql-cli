package generate

import (
	"fmt"
	"github.com/FINTLabs/fint-graphql-cli/common/config"
	"github.com/FINTLabs/fint-graphql-cli/common/document"
	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/FINTLabs/fint-graphql-cli/common/parser"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strings"
)

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")
	owner := c.GlobalString("owner")
	repo := c.GlobalString("repo")
	filename := c.GlobalString("filename")


	setupGraphQlSchemaDirStructure()
	generateGraphQlSchema(owner, repo, tag, filename, force)
	generateGraphQlQueryResolver(owner, repo, tag, filename, force)
	generateGraphQlService(owner, repo, tag, filename, force)
	generateGraphQlResolver(owner, repo, tag, filename, force)
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

func writeSchema(schema string, content []byte) error {
	path := fmt.Sprintf("%s/schema", config.GRAPHQL_BASE_PATH)
	return writeFile(path, strings.ToLower(schema) + ".graphqls", []byte(content))
}

func writeQueryResolver(className string, content []byte) error {
	path := fmt.Sprintf("%s/model/%s", config.GRAPHQL_BASE_PATH, strings.ToLower(className))
	return writeFile(path, fmt.Sprintf("%sQueryResolver.java", className)																																							, []byte(content))
}

func writeService(className string, content []byte) error {
	path := fmt.Sprintf("%s/model/%s", config.GRAPHQL_BASE_PATH, strings.ToLower(className))
	return writeFile(path, fmt.Sprintf("%sService.java", className)																																							, []byte(content))
}

func writeResolver(className string, content []byte) error {
	path := fmt.Sprintf("%s/model/%s", config.GRAPHQL_BASE_PATH, strings.ToLower(className))
	return writeFile(path, fmt.Sprintf("%sResolver.java", className)																																							, []byte(content))
}

func generateGraphQlSchema(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating GraphQL Schema")


	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	for _, c := range classes {
		if !c.Abstract && includePackage(c.Package) {
			fmt.Printf("  > Creating schema: %s.graphqls\n", c.Name)
			schema := GetGraphQlSchema(c)
			err := writeSchema(c.Name, []byte(schema))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}




}

func generateGraphQlQueryResolver(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating GraphQL Query Resolver")

	//setupGraphQlSchemaDirStructure()

	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	for _, c := range classes {
		if !c.Abstract && c.Stereotype == "hovedklasse" && includePackage(c.Package) {
			fmt.Printf("  > Creating query resolver: %s.java\n", c.Name)
			class := GetGraphQlQueryReolver(c)
			err := writeQueryResolver(c.Name, []byte(class))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}
}

func generateGraphQlService(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating GraphQL Service")

	//setupGraphQlSchemaDirStructure()

	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	for _, c := range classes {
		if !c.Abstract && c.Stereotype == "hovedklasse" && includePackage(c.Package) {
			fmt.Printf("  > Creating service: %s.java\n", c.Name)
			class := GetGraphQlService(c)
			err := writeService(c.Name, []byte(class))
			if err != nil {
				fmt.Printf("Unable to write file: %s", err)
			}
		}
	}
}

func generateGraphQlResolver(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating GraphQL Resolver")

	//setupGraphQlSchemaDirStructure()

	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	for _, c := range classes {
		if !c.Abstract && c.Stereotype == "hovedklasse" && includePackage(c.Package) {
			fmt.Printf("  > Creating service: %s.java\n", c.Name)
			class := GetGraphQlResolver(c)
			err := writeResolver(c.Name, []byte(class))
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

