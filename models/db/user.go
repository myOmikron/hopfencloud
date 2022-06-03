package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type User struct {
	utilitymodels.Common
	AuthID       uint
	AuthKey      string
	MailVerified bool
}
