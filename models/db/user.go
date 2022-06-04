package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type AccountEmailVerification struct {
	utilitymodels.Common
	AccountID uint
	Account   Account
	Email     string
	Token     string
}

type Account struct {
	utilitymodels.Common
	AuthID  uint
	AuthKey string
}
