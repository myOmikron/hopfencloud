package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/myOmikron/hopfencloud/server"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("hopfencloud", "", &argparse.ParserConfig{})
	startParser := parser.AddCommand("start", "Starts the server", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})
	defaultConf := startParser.String("", "config-path", &argparse.Option{
		Default: "/etc/hopfencloud/config.toml",
	})

	reloadParser := parser.AddCommand("reload", "Sends SIGHUP the the given process", nil)
	reloadID := reloadParser.Int("", "process-id", &argparse.Option{
		Required:   true,
		Positional: true,
		Help:       "Process ID",
		Validate: func(arg string) error {
			if i, _ := strconv.Atoi(arg); i <= 0 {
				return errors.New("int must be > 0")
			}
			return nil
		},
	})

	stopParser := parser.AddCommand("stop", "Sends SIGINT to the given process", nil)
	stopID := stopParser.Int("", "process-id", &argparse.Option{
		Required:   true,
		Positional: true,
		Help:       "Process ID",
		Validate: func(arg string) error {
			if i, _ := strconv.Atoi(arg); i <= 0 {
				return errors.New("int must be > 0")
			}
			return nil
		},
	})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case startParser.Invoked:
		server.StartServer(*defaultConf, false)
	case reloadParser.Invoked:
		if process, err := os.FindProcess(*reloadID); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			if err := process.Signal(syscall.SIGHUP); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	case stopParser.Invoked:
		if process, err := os.FindProcess(*stopID); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			if err := process.Signal(syscall.SIGINT); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}
