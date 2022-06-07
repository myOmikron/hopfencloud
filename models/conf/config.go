package conf

import (
	"errors"
	"os"
	"path"
)

type Database struct {
	Driver   string
	Name     string
	Port     uint
	Host     string
	User     string
	Password string
}

type AllowedHost struct {
	Host  string
	Https bool
}

type Server struct {
	CLISockPath             string
	ListenAddress           string
	AllowedHosts            []*AllowedHost
	UseForwardedProtoHeader bool
}

type Files struct {
	DataPath string
}

type Config struct {
	Database Database
	Files    Files
	Server   Server
}

func (c *Config) CheckConfig() error {
	if c.Files.DataPath == "" {
		return errors.New("invalid Files.DataPath setting")
	}

	if stat, err := os.Stat(c.Files.DataPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(c.Files.DataPath, 0700); err != nil {
				return err
			}
		}
	} else {
		if !stat.IsDir() {
			return errors.New("data directory is a file")
		}

		if _, err := os.OpenFile(path.Join(c.Files.DataPath, "test"), os.O_CREATE, 0600); err != nil {
			return err
		}

		if err := os.Remove(path.Join(c.Files.DataPath, "test")); err != nil {
			return err
		}
	}

	return nil
}
