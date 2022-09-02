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

	settingsParser := parser.AddCommand("settings", "Configure and show server settings", nil)
	settingsShowParser := settingsParser.AddCommand(
		"show",
		"Show the current settings",
		&argparse.ParserConfig{
			DisableDefaultShowHelp: true,
		},
	)

	if err := parser.Parse(nil); err != nil {
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	switch {
	case createAdminUserParser.Invoked:
		cli.CreateAdminUser(*sockPath, createAdminUserArgs)
	case settingsParser.Invoked:
		settingsParser.PrintHelp()
	case settingsShowParser.Invoked:
		cli.SettingsShow(*sockPath, &cli.SettingsCLI{})
	default:
		parser.PrintHelp()
	}
}
