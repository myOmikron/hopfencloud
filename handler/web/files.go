package web

import (
	"strconv"
	"time"

	"github.com/myOmikron/hopfencloud/modules/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FileQuery struct {
	Dir *uint `query:"dir"`
}

type File struct {
	ID           uint
	Name         string
	Size         string
	LastModified string
	IsDirectory  bool
}

type FileData struct {
	PageTitle          string
	CurrentDirectoryID string
	Files              []File
}

func (w *Wrapper) Files(c echo.Context) error {
	var query FileQuery

	if err := c.Bind(&query); err != nil {
		//TODO: Display error message
		return c.String(500, "Internal server error")
	}

	account, err := utils.GetAccount(c, w.DB)
	if err != nil {
		//TODO: Display error message
		return c.String(500, "Internal server error")
	}

	//TODO: As this only gets slow, if a user has very many files / directories
	// in the current directory (> 2000), there's currently no limit / pagination
	//
	// Add this to Preload
	//	func(db *gorm.DB) *gorm.DB {
	//		return db.Offset(int(query.Offset)).Limit(1000)
	//	},
	w.DB.Preload(
		"Files",
		map[string]interface{}{
			"parent_id": query.Dir,
		},
	).Preload(
		"InternalShares",
		func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "target_id", "target_type")
		},
	).Find(&account)

	files := make([]File, len(account.Files))
	for _, file := range account.Files {
		files = append(files, File{
			ID:           file.ID,
			Name:         file.Name,
			Size:         utils.ByteCountDecimal(file.Size),
			LastModified: file.FileUpdatedAt.Format(time.RFC1123Z),
			IsDirectory:  file.IsDirectory,
		})
	}

	var current string
	if query.Dir != nil {
		current = strconv.Itoa(int(*query.Dir))
	}
	return c.Render(200, "files", &FileData{
		PageTitle:          "Files - " + w.Settings.SiteName,
		Files:              files,
		CurrentDirectoryID: current,
	})
}
