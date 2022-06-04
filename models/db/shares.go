package db

import (
	"time"

	"github.com/myOmikron/echotools/utilitymodels"
)

type ExternalShare struct {
	utilitymodels.Common
	Name     string `gorm:"size:256"`
	Token    string `gorm:"size:256"`
	Read     bool
	Upload   bool
	Expired  *time.Time
	Password *string `gorm:"size:256"`
}

type InternalShare struct {
	utilitymodels.Common
	Read          bool
	Upload        bool
	Expired       *time.Time
	VirtualUserID uint
	VirtualUser   VirtualUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
