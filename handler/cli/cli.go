package cli

import (
	"text/template"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"

	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type CLI struct {
	Config             *conf.Config
	DB                 *gorm.DB
	MailTemplates      *template.Template
	ReloadSettingsFunc func()
	Settings           *db.Settings
	WorkerPool         worker.Pool
}
