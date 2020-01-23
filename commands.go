package main

import (
	"fmt"
	"os"

	"github.com/FINTLabs/fint-graphql-cli/classes"
	"github.com/FINTLabs/fint-graphql-cli/generate"
	"github.com/FINTLabs/fint-graphql-cli/packages"

	"github.com/FINTLabs/fint-graphql-cli/branches"
	"github.com/FINTLabs/fint-graphql-cli/tags"
	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "GITHUB_OWNER",
		Name:   "owner",
		Value:  "FINTLabs",
		Usage:  "Git repository containing model",
	},
	cli.StringFlag{
		EnvVar: "GITHUB_PROJECT",
		Name:   "repo",
		Value:  "fint-informasjonsmodell",
		Usage:  "Git repository containing model",
	},
	cli.StringFlag{
		EnvVar: "MODEL_FILENAME",
		Name:   "filename",
		Value:  "FINT-informasjonsmodell.xml",
		Usage:  "File name containing information model",
	},
	cli.StringFlag{
		EnvVar: "",
		Name:   "tag, t",
		Value:  "latest",
		Usage:  "the tag (version) of the model to generate",
	},
	cli.BoolFlag{
		EnvVar: "",
		Name:   "force, f",
		Usage:  "force downloading XMI for GitHub.",
	},
}

var Commands = []cli.Command{

	{
		Name:   "printClasses",
		Usage:  "list classes",
		Action: classes.CmdPrintClasses,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "generate",
		Usage:  "generates GraphQL schema, query resolvers, resolvers and services.",
		Action: generate.CmdGenerate,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "exclude, x",
				Usage: "Classes to exclude from the generated schema and classes.",
			},
		},
	},
	{
		Name:   "listPackages",
		Usage:  "list Packages",
		Action: packages.CmdListPackages,
		Flags:  []cli.Flag{},
	},
	/*
		{
			Name:   "listNamespaces",
			Usage:  "list CS namespaces",
			Action: namespaces.CmdListNamespaces,
			Flags:  []cli.Flag{},
		},
	*/
	{
		Name:   "listTags",
		Usage:  "list tags",
		Action: tags.CmdListTags,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listBranches",
		Usage:  "list branches",
		Action: branches.CmdListBranches,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
