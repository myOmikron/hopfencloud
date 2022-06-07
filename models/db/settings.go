package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Settings struct {
	utilitymodels.Common
	//General
	SiteName      string `gorm:"size:256" json:"site_name"`
	PublicAddress string `gorm:"size:256" json:"public_address"` // Used for links pointing to this server, e.g. confirmation mail link generation

	//Authentication
	LocalRegistrationDisabled bool `json:"local_registration_disabled"`

	//SMTP related settings
	SMTPFrom     string `gorm:"size:256" json:"smtp_from"`
	SMTPHost     string `gorm:"size:256" json:"smtp_host"`
	SMTPPort     uint16 `json:"smtp_port"`
	SMTPUser     string `gorm:"size:256" json:"smtp_user"`
	SMTPPassword string `gorm:"size:256" json:"smtp_password"`
}
