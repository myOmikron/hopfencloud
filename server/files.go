package server

import (
	"os"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/utils"

	"gorm.io/gorm"
)

func initializeDirStructure(dbase *gorm.DB, config *conf.Config) error {
	// User
	if err := os.MkdirAll(utils.GetUserPath(config), 0700); err != nil {
		return err
	}

	accounts := make([]db.Account, 0)
	dbase.Find(&accounts)
	for _, account := range accounts {
		if err := os.MkdirAll(utils.GetUserCurrentPath(account.ID, config), 0700); err != nil {
			return err
		}

		if err := os.MkdirAll(utils.GetUserVersionsPath(account.ID, config), 0700); err != nil {
			return err
		}
	}

	// Groups
	if err := os.MkdirAll(utils.GetGroupPath(config), 0700); err != nil {
		return err
	}

	groups := make([]db.Group, 0)
	dbase.Find(&groups)
	for _, group := range groups {
		if err := os.MkdirAll(utils.GetGroupCurrentPath(group.ID, config), 0700); err != nil {
			return err
		}

		if err := os.MkdirAll(utils.GetGroupVersionsPath(group.ID, config), 0700); err != nil {
			return err
		}
	}

	return nil
}
