package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type AccountEmailVerification struct {
	utilitymodels.Common
	AccountID uint
	Account   Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email     string  `gorm:"size:256"`
	Token     string  `gorm:"size:256"`
}

type Account struct {
	utilitymodels.Common
	AuthID  uint
	AuthKey string `gorm:"size:256"`
}
