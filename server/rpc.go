package server

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"text/template"

	"github.com/myOmikron/hopfencloud/handler/cli"
	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"github.com/myOmikron/echotools/worker"
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
		if err := os.Chmod(config.Server.CLISockPath, 0777); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Start listening on " + config.Server.CLISockPath)
		if err := http.Serve(*sock, nil); !errors.Is(err, net.ErrClosed) {
			logger.Error(err.Error())
		}
	}
}
