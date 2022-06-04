package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Group struct {
	utilitymodels.Common
	Name           string      `gorm:"size:256"`
	VirtualUser    VirtualUser `gorm:"polymorphic:Model;"`
	LinkedAccounts []Account   `gorm:"many2many:group__linked_accounts"`
}
