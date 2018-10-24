package types

import (
	"fmt"
	"strings"
)

const LIST_TEMPLATE = "List<%s>"

func getType(list bool, t string) string {
	if list {
		return fmt.Sprintf(LIST_TEMPLATE, t)
	}
	return t
}

func GetComponentName(pkg string) string {
	return strings.ToLower(strings.Split(pkg, ".")[3])
}
