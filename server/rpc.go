package server

import (
	"errors"
	"github.com/myOmikron/echotools/worker"
	"net"
	"net/http"
	"net/rpc"
	"text/template"

	"github.com/myOmikron/hopfencloud/handler/cli"
	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"gorm.io/gorm"
)

func initializeRPC(
	sock *net.Listener,
	db *gorm.DB,
	config *conf.Config,
	workerPool worker.Pool,
	mailTemplates *template.Template,
	reloadSettingsFunc func(),
	settings *db.Settings,
	isReloading bool,
) {
	if !isReloading {
		handler := cli.CLI{
			Config:             config,
			DB:                 db,
			MailTemplates:      mailTemplates,
			ReloadSettingsFunc: reloadSettingsFunc,
			Settings:           settings,
			WorkerPool:         workerPool,
		}
		if err := rpc.Register(&handler); err != nil {
			logger.Error(err.Error())
		}
		rpc.HandleHTTP()
	}

	var err error
	if *sock, err = net.Listen("unix", config.Server.CLISockPath); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("Start listening on " + config.Server.CLISockPath)
		if err := http.Serve(*sock, nil); !errors.Is(err, net.ErrClosed) {
			logger.Error(err.Error())
		}
	}
}
