package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type UserMailConfirmation struct {
	utilitymodels.Common
	UserID uint
	User   User
	Mail   string
	Token  string
}

type User struct {
	utilitymodels.Common
	AuthID  uint
	AuthKey string
}
