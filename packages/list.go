package packages

import (
	"github.com/FINTLabs/fint-graphql-cli/common/types"
	"github.com/FINTLabs/fint-graphql-cli/common/utils"
)

func DistinctPackageList(classes []*types.Class) []string {

	var p []string
	for _, c := range classes {
		p = append(p, c.Package)
	}

	return utils.Distinct(p)
}
