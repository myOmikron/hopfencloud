package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Settings struct {
	utilitymodels.Common
	//General
	SiteName      string `gorm:"size:256"`
	PublicAddress string `gorm:"size:256"` // Used for links pointing to this server, e.g. confirmation mail link generation

	//Authentication
	LocalRegistrationDisabled bool

	//SMTP related settings
	SMTPFrom     string `gorm:"size:256"`
	SMTPHost     string `gorm:"size:256"`
	SMTPPort     uint16
	SMTPUser     string `gorm:"size:256"`
	SMTPPassword string `gorm:"size:256"`
}
