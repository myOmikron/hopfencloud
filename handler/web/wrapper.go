package web

import (
	"text/template"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"

	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type Wrapper struct {
	DB                 *gorm.DB
	Config             *conf.Config
	WorkerPool         worker.Pool
	MailTemplates      *template.Template
	Settings           *db.Settings
	SettingsReloadFunc func()
}
