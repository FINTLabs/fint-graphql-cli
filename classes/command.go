package classes

import (
	"fmt"

	"github.com/FINTLabs/fint-graphql-cli/common/github"
	"github.com/FINTLabs/fint-graphql-cli/common/parser"
	"github.com/FINTLabs/fint-graphql-cli/common/types"
	"github.com/codegangsta/cli"
)

func CmdPrintClasses(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}

	clazzes, _, _, _ := parser.GetClasses(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), c.GlobalBool("force"))

	for _, c := range clazzes {
		dumpClass(c)
	}

}

func dumpClass(c *types.Class) {
	dep := ""
	if c.Deprecated {
		dep = "<<DEPRECATED>>"
	}
	fmt.Printf("Class: %s %s\n", c.Name, dep)
	fmt.Printf("  Abstract: %t\n", c.Abstract)
	fmt.Printf("  Identifiable: %t\n", c.Identifiable)
	fmt.Printf("  Resource: %t\n", c.Resource)
	fmt.Printf("  Package: %s\n", c.Package)
	fmt.Printf("  Namespace: %s\n", c.Namespace)
	//fmt.Printf("  DocumentationUrl: %s\n", c.DocumentationUrl)
	if len(c.Extends) > 0 {
		fmt.Printf("  Extends: %s\n", c.Extends)
		fmt.Printf("  ExtendsResource: %t\n", c.ExtendsResource)
	}
	fmt.Println("  Imports:")
	for _, i := range c.Imports {
		fmt.Printf("    - %s\n", i)
	}
	fmt.Println("  Using:")
	for _, u := range c.Using {
		fmt.Printf("    - %s\n", u)
	}

	if len(c.Identifiers) > 0 {
		fmt.Println("  Identifiers:")
		for _, i := range c.Identifiers {
			fmt.Printf("    - %s\n", i.Name)
		}
	}

	if len(c.Attributes) > 0 {
		fmt.Println("  Attributes:")
		for _, a := range c.Attributes {
			dep := ""
			if a.Deprecated {
				dep = "<<DEPRECATED>>"
			}
			if a.List {
				fmt.Printf("    - %s: List<%s> %s\n", a.Name, a.Type, dep)
			} else {
				fmt.Printf("    - %s: %s %s\n", a.Name, a.Type, dep)
			}
		}
	}

	if len(c.InheritedAttributes) > 0 {
		fmt.Println("  Inherited Attributes:")
		for _, a := range c.InheritedAttributes {
			dep := ""
			if a.Deprecated {
				dep = "<<DEPRECATED>>"
			}
			if a.List {
				fmt.Printf("    - %s: List<%s> [%s] %s\n", a.Name, a.Type, a.Owner, dep)
			} else {
				fmt.Printf("    - %s: %s [%s] %s\n", a.Name, a.Type, a.Owner, dep)
			}
		}
	}

	if len(c.Relations) > 0 {
		fmt.Println("  Relations:")
		for _, r := range c.Relations {
			s := ""
			if r.Deprecated {
				s = "<<DEPRECATED>>"
			}
			m := ""
			if r.Optional && r.List {
				m = "0..*"
			} else if r.List {
				m = "1..*"
			} else {
				m = "1"
			}

			fmt.Printf("    - %s: %s[%s] <<%s>> %s\n", r.Name, r.Target, m, r.Stereotype, s)
		}
	}

	if len(c.Resources) > 0 {
		fmt.Println("  Resources:")
		for _, a := range c.Resources {
			if a.List {
				fmt.Printf("    - %s: List<%s>\n", a.Name, a.Type)
			} else {
				fmt.Printf("    - %s: %s\n", a.Name, a.Type)
			}
		}
	}
}
