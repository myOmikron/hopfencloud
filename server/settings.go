package server

import (
	"time"

	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"gorm.io/gorm"
)

func getSettingsReloadFunc(settings *db.Settings, database *gorm.DB) func() {
	database.First(&settings)

	c := make(chan bool)
	go func() {
		for {
			select {
			case <-c:
				database.First(&settings)
				logger.Info("Reloaded settings")
			default: // Don't block
			}
			<-time.After(time.Second)
		}
	}()

	return func() {
		c <- true
	}
}
