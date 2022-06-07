package utils

import (
	"errors"

	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/gorm"
)

func GetAccount(c echo.Context, dbase *gorm.DB) (*db.Account, error) {
	sessionContext, err := middleware.GetSessionContext(c)
	if err != nil {
		return nil, err
	}

	var authKey string
	var authID uint

	user := sessionContext.GetUser()
	switch user := user.(type) {
	case *utilitymodels.LocalUser:
		authKey = "local"
		authID = user.ID
	case *utilitymodels.LDAPUser:
		authKey = "ldap"
		authID = user.ID
	}

	account := db.Account{}
	var count int64
	dbase.Find(&account, "auth_id = ? AND auth_key = ?", authID, authKey).Count(&count)
	if count != 1 {
		logger.Info("Could not determine account")
		return nil, errors.New("could not determine account")
	}

	return &account, nil

}
