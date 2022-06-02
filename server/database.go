package server

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//cleanupDatabase is a task that will run periodically every 5 minutes to do
// some cleanup in the database:
// - Remove expired sessions
// - Removed sessions when the linked user doesn't exist anymore
func cleanupDatabase(db *gorm.DB) {
	var start time.Time
	for {
		start = time.Now()

		// Delete expired sessions
		db.Delete(&utilitymodels.Session{}, "valid_until < ?", start)

		// Delete sessions if the linked user does not exist anymore
		sessions := make([]utilitymodels.Session, 0)
		db.Find(&sessions)

		// It will be cheaper to reallocate those
		localMap := make(map[uint]bool)
		ldapMap := make(map[uint]bool)

		toDeleteSessions := make([]uint, 0)

		var count int64
		for _, session := range sessions {
			switch session.AuthKey {
			case "local":
				if value, exists := localMap[session.AuthID]; exists {
					if !value {
						toDeleteSessions = append(toDeleteSessions, session.ID)
					}
				} else {
					db.Find(&utilitymodels.LocalUser{}, "ID = ?", session.AuthID).Count(&count)
					if count == 0 {
						localMap[session.AuthID] = false
						toDeleteSessions = append(toDeleteSessions, session.ID)
					} else {
						localMap[session.AuthID] = true
					}
				}
			case "ldap":
				if value, exists := ldapMap[session.AuthID]; exists {
					if !value {
						toDeleteSessions = append(toDeleteSessions, session.ID)
					}
				} else {
					db.Find(&utilitymodels.LDAPUser{}, "ID = ?", session.AuthID).Count(&count)
					if count == 0 {
						ldapMap[session.AuthID] = false
						toDeleteSessions = append(toDeleteSessions, session.ID)
					} else {
						ldapMap[session.AuthID] = true
					}
				}
			}
		}

		if len(toDeleteSessions) > 0 {
			db.Delete(&utilitymodels.Session{}, toDeleteSessions)
		}

		// Cleanup
		logger.Infof("Cleanup took %s", time.Now().Sub(start).String())

		// Sleeping
		time.Sleep(time.Minute * 5)
	}
}

func initializeDatabase(config *conf.Config) *gorm.DB {
	// Database
	var driver gorm.Dialector
	switch config.Database.Driver {
	case "sqlite":
		driver = sqlite.Open(config.Database.Name)
	case "mysql":
		mysqlConf := mysqlDriver.NewConfig()
		mysqlConf.Net = fmt.Sprintf("tcp(%s)", net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))))
		mysqlConf.DBName = config.Database.Name
		mysqlConf.User = config.Database.User
		mysqlConf.Passwd = config.Database.Password
		mysqlConf.ParseTime = true
		mysqlConf.Params = map[string]string{
			"charset": "utf8mb4",
		}
		driver = mysql.Open(mysqlConf.FormatDSN())
	case "postgresql":
		dsn := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(config.Database.User, config.Database.Password),
			Host:   net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))),
			Path:   config.Database.Name,
		}
		driver = postgres.Open(dsn.String())
	}

	dbase := database.Initialize(
		driver,
	)

	return dbase
}
