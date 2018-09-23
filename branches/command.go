package branches

import (
	"fmt"

	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/codegangsta/cli"
)

func CmdListBranches(c *cli.Context) {
	for _, b := range github.GetBranchList(c.GlobalString("owner"), c.GlobalString("repo")) {
		fmt.Println(b)
	}
}
