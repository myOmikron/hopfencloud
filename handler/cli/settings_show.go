package cli

import (
	"github.com/myOmikron/hopfencloud/models/db"
)

type SettingsShowRequest struct {
}

type SettingsShowResult struct {
	ErrorMessage *string
	Settings     *db.Settings
}

func (c *CLI) SettingsShow(req SettingsShowRequest, res *SettingsShowResult) error {
	res.Settings = c.Settings
	return nil
}
