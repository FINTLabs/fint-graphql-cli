package namespaces

import (
	"github.com/FINTLabs/fint-graphql-cli/common/parser"
	"github.com/FINTLabs/fint-graphql-cli/common/utils"
)

func DistinctNamespaceList(owner string, repo string, tag string, filename string, force bool) []string {
	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	var p []string
	for _, c := range classes {
		p = append(p, c.Namespace)
	}

	return utils.Distinct(p)
}
