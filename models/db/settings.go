package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Settings struct {
	utilitymodels.Common
	SiteName                  string
	LocalRegistrationDisabled bool

	//SMTP related settings
	SMTPFrom     string
	SMTPHost     string
	SMTPPort     uint
	SMTPUser     string
	SMTPPassword string
}
