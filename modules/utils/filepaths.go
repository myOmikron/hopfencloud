package utils

import (
	"path/filepath"
	"strconv"

	"github.com/myOmikron/hopfencloud/models/conf"
)

func GetUserPath(config *conf.Config) string {
	return filepath.Join(config.Files.DataPath, "users")
}

func GetUserCurrentPath(userID uint, config *conf.Config) string {
	return filepath.Join(GetUserPath(config), strconv.Itoa(int(userID)), "current")
}

func GetUserVersionsPath(userID uint, config *conf.Config) string {
	return filepath.Join(GetUserPath(config), strconv.Itoa(int(userID)), "versions")
}

func GetGroupPath(config *conf.Config) string {
	return filepath.Join(config.Files.DataPath, "groups")
}

func GetGroupCurrentPath(groupID uint, config *conf.Config) string {
	return filepath.Join(GetGroupPath(config), strconv.Itoa(int(groupID)), "current")
}

func GetGroupVersionsPath(groupID uint, config *conf.Config) string {
	return filepath.Join(GetGroupPath(config), strconv.Itoa(int(groupID)), "versions")
}
