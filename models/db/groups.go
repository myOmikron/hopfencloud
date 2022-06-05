package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Group struct {
	utilitymodels.Common
	Name           string          `gorm:"size:256"`
	Files          []File          `gorm:"polymorphic:Owner;"`
	InternalShares []InternalShare `gorm:"polymorphic:Target;"`
	LinkedAccounts []Account       `gorm:"many2many:group__linked_accounts"`
}
