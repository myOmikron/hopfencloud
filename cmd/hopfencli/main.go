package main

import (
	"os"

	"github.com/myOmikron/hopfencloud/cli"

	"github.com/hellflame/argparse"
	"github.com/myOmikron/echotools/color"
)

func main() {
	parser := argparse.NewParser("hopfencli", "CLI tools for hopfencloud", nil)

	sockPath := parser.String("", "sock-path", &argparse.Option{
		Inheritable: true,
		Default:     "/run/hopfencloud/cli.sock",
		Help:        "Alternative path to cli socket. Default to /run/hopfencloud/cli.sock",
	})

	createAdminUserArgs, createAdminUserParser := cli.RegisterCreateAdminUser(parser)

	if err := parser.Parse(nil); err != nil {
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	switch {
	case createAdminUserParser.Invoked:
		cli.CreateAdminUser(*sockPath, createAdminUserArgs)
	default:
		parser.PrintHelp()
	}
}
