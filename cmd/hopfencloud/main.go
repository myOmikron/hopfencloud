package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"hopfencloud/server"
	"os"
)

func main() {
	parser := argparse.NewParser("hopfencloud", "", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	defaultConf := parser.String("", "config-path", &argparse.Option{
		Default: "/etc/hopfencloud/config.toml",
	})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	server.StartServer(*defaultConf)
}
