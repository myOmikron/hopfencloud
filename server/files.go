package server

import (
	"os"
	"time"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"github.com/myOmikron/hopfencloud/modules/utils"

	"gorm.io/gorm"
)

type result struct {
	FileID uint
	Total  uint
}

func cleanupVersionedFiles(dbase *gorm.DB, config *conf.Config, settings *db.Settings) {
	var start time.Time
	for {
		start = time.Now()

		if !settings.VersioningKeepAll {
			results := make([]result, 0)
			fileVersions := make([]db.FileVersion, 0)
			toDeleteFileVersions := make([]db.FileVersion, 0)
			var count int64
			var file db.File

			// Keep Last
			for {
				dbase.Table("file__file_versions").
					Limit(1000).
					Select("file_id, COUNT(fv.id) AS total").
					Joins("join file_versions fv on file__file_versions.file_version_id = fv.id").
					Group("file_id").
					Having("count(*) > ?", settings.VersioningKeepLast).
					Scan(&results).
					Count(&count)

				for _, res := range results {
					file.ID = res.FileID
					dbase.Model(&file).
						Select("id, path").
						Order("file_relevant_until asc").
						Limit(int(res.Total - settings.VersioningKeepLast)).
						Association("FileVersions").
						Find(&fileVersions)

					dbase.Model(&file).Association("FileVersions").Delete(&fileVersions)
					for _, version := range fileVersions {
						toDeleteFileVersions = append(toDeleteFileVersions, version)
					}
				}

				if len(toDeleteFileVersions) > 0 {
					dbase.Delete(&toDeleteFileVersions)
					for _, fileVersion := range toDeleteFileVersions {
						if err := os.Remove(fileVersion.Path); err != nil {
							logger.Error("File cleanup: " + err.Error())
						}
					}
				}

				if count < 1000 {
					break
				}
			}

			// Keep hourly
			if settings.VersioningKeepHourly != nil {

				for {
					dbase.Table("file__file_versions").
						Limit(1000).
						Select("file_id, COUNT(fv.id) AS total").
						Joins("join file_versions fv on file__file_versions.file_version_id = fv.id").
						Group("file_id").
						Having("count(*) > ?", *settings.VersioningKeepHourly).
						Scan(&results).
						Count(&count)

					for _, res := range results {
						file.ID = res.FileID
						dbase.Model(&file).
							Select("id, path").
							Order("file_relevant_until asc").
							Limit(int(res.Total - settings.VersioningKeepLast)).
							Association("FileVersions").
							Find(&fileVersions)

						dbase.Model(&file).Association("FileVersions").Delete(&fileVersions)
						for _, version := range fileVersions {
							toDeleteFileVersions = append(toDeleteFileVersions, version)
						}
					}

					if len(toDeleteFileVersions) > 0 {
						dbase.Delete(&toDeleteFileVersions)
						for _, fileVersion := range toDeleteFileVersions {
							if err := os.Remove(fileVersion.Path); err != nil {
								logger.Error("File cleanup: " + err.Error())
							}
						}
					}

					if count < 1000 {
						break
					}
				}
			}
		}

		logger.Infof("File cleanup took %s", time.Now().Sub(start))
		time.Sleep(time.Minute * 10)
	}
}

func initializeDirStructure(dbase *gorm.DB, config *conf.Config) error {
	// User
	if err := os.MkdirAll(utils.GetUserPath(config), 0700); err != nil {
		return err
	}

	accounts := make([]db.Account, 0)
	dbase.Find(&accounts)
	for _, account := range accounts {
		if err := os.MkdirAll(utils.GetUserCurrentPath(account.ID, config), 0700); err != nil {
			return err
		}

		if err := os.MkdirAll(utils.GetUserVersionsPath(account.ID, config), 0700); err != nil {
			return err
		}
	}

	// Groups
	if err := os.MkdirAll(utils.GetGroupPath(config), 0700); err != nil {
		return err
	}

	groups := make([]db.Group, 0)
	dbase.Find(&groups)
	for _, group := range groups {
		if err := os.MkdirAll(utils.GetGroupCurrentPath(group.ID, config), 0700); err != nil {
			return err
		}

		if err := os.MkdirAll(utils.GetGroupVersionsPath(group.ID, config), 0700); err != nil {
			return err
		}
	}

	return nil
}
