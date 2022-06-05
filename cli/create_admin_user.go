package cli

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"net/rpc"
	"os"
	"strings"
	"syscall"

	"github.com/myOmikron/hopfencloud/handler/cli"

	"github.com/hellflame/argparse"
	"github.com/myOmikron/echotools/color"
)

type CreateAdminUserCLI struct {
}

func RegisterCreateAdminUser(parser *argparse.Parser) (*CreateAdminUserCLI, *argparse.Parser) {
	var args CreateAdminUserCLI

	createAdminUserParser := parser.AddCommand(
		"create-admin-user",
		"Create a local admin user",
		&argparse.ParserConfig{
			DisableDefaultShowHelp: true,
		},
	)

	return &args, createAdminUserParser
}

func CreateAdminUser(sockPath string, args *CreateAdminUserCLI) {
	if conn, err := rpc.DialHTTP("unix", sockPath); err != nil {
		color.Println(color.RED, err.Error())
	} else {
		req := cli.CreateAdminUserRequest{}
		res := cli.CreateAdminUserResponse{}

		var reader = bufio.NewReader(os.Stdin)

		color.Print(color.BLUE, "Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		color.Print(color.BLUE, "Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		color.Print(color.BLUE, "Password: ")
		bytePassword, _ := term.ReadPassword(syscall.Stdin)
		fmt.Println()

		req.Email = email
		req.Username = username
		req.Password = string(bytePassword)

		if err := conn.Call("CLI.CreateAdminUser", req, &res); err != nil {
			color.Println(color.RED, "Error returned:")
			fmt.Println(err.Error())
			if res.ErrorMessage != nil {
				fmt.Println("Additional message:", *res.ErrorMessage)
			}
			os.Exit(1)
		}

		color.Println(color.PURPLE, "Admin created successfully!")
	}
}
