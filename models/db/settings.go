package db

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type Settings struct {
	utilitymodels.Common
	SiteName                  string
	LocalRegistrationDisabled bool
}
