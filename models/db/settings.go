package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Settings struct {
	utilitymodels.Common
	//General
	SiteName      string
	PublicAddress string // Used for links pointing to this server, e.g. confirmation mail link generation

	//Authentication
	LocalRegistrationDisabled bool

	//SMTP related settings
	SMTPFrom     string
	SMTPHost     string
	SMTPPort     uint
	SMTPUser     string
	SMTPPassword string
}
