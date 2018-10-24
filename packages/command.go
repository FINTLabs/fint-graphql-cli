package packages

import (
	"fmt"

	"github.com/FINTLabs/fint-graphql-cli/common/config"
	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/FINTLabs/fint-graphql-cli/common/parser"
	"github.com/codegangsta/cli"
)

func CmdListPackages(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}

	classes, _, _, _ := parser.GetClasses(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), c.GlobalBool("force"))

	for _, p := range DistinctPackageList(classes) {
		fmt.Println(p)
	}
}
