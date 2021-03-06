package namespaces

import (
	"fmt"

	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/urfave/cli"
)

func CmdListNamespaces(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}

	for _, p := range DistinctNamespaceList(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), c.GlobalBool("force")) {
		fmt.Println(p)
	}
}
