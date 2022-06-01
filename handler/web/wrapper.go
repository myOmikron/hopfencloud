package web

import (
	"github.com/myOmikron/hopfencloud/models/conf"

	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type Wrapper struct {
	DB         *gorm.DB
	Config     *conf.Config
	WorkerPool worker.Pool
}
