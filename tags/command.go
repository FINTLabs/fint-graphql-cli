package tags

import (
	"fmt"

	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/codegangsta/cli"
)

func CmdListTags(c *cli.Context) {
	for _, t := range github.GetTagList(c.GlobalString("owner"), c.GlobalString("repo")) {
		fmt.Println(t)
	}
}