package cli

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/rpc"
	"os"

	"github.com/myOmikron/hopfencloud/handler/cli"

	"github.com/myOmikron/echotools/color"
)

type SettingsCLI struct {
}

var json = jsoniter.Config{
	EscapeHTML:    true,
	CaseSensitive: true,
}.Froze()

func SettingsShow(sockPath string, args *SettingsCLI) {
	if conn, err := rpc.DialHTTP("unix", sockPath); err != nil {
		color.Println(color.RED, err.Error())
	} else {
		var res cli.SettingsShowResult

		if err := conn.Call("CLI.SettingsShow", &cli.SettingsShowRequest{}, &res); err != nil {
			color.Println(color.RED, "Error returned:")
			fmt.Println(err.Error())
			if res.ErrorMessage != nil {
				fmt.Println("Additional message:", *res.ErrorMessage)
			}
			os.Exit(1)
		}

		color.Println(color.BLUE, "[SMTP]")
		fmt.Print("Host: ")
		color.Println(color.CYAN, res.Settings.SMTPHost)
		fmt.Print("Port: ")
		color.Println(color.CYAN, fmt.Sprint(res.Settings.SMTPPort))
		fmt.Print("From: ")
		color.Println(color.CYAN, res.Settings.SMTPFrom)
		fmt.Print("User: ")
		color.Println(color.CYAN, res.Settings.SMTPUser)
		fmt.Print("Password: ")
		color.Println(color.CYAN, res.Settings.SMTPPassword)

		fmt.Println()

	}
}
