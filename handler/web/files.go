package web

import (
	"github.com/myOmikron/hopfencloud/models/db"
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

type PathEntry struct {
	ID       uint
	Path     string
	ParentID *uint
	Name     string
}

type FileData struct {
	PageTitle          string
	IsAdmin            bool
	CurrentDirectoryID int
	ParentDirectoryID  int
	Files              []File
	Path               []PathEntry
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
	).
		Preload(
			"InternalShares",
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "target_id", "target_type")
			},
		).Find(&account)

	directories := make([]File, 0)
	files := make([]File, 0)
	for _, file := range account.Files {
		if file.IsDirectory {
			directories = append(directories, File{
				ID:          file.ID,
				Name:        file.Name,
				IsDirectory: true,
			})
		} else {
			files = append(files, File{
				ID:           file.ID,
				Name:         file.Name,
				Size:         utils.ByteCountDecimal(file.Size),
				LastModified: file.FileUpdatedAt.Format(time.RFC1123Z),
				IsDirectory:  file.IsDirectory,
			})
		}
	}

	var current int
	var parent int
	path := make([]PathEntry, 0)

	if query.Dir != nil {
		current = int(*query.Dir)
		var currentDirectory db.File
		if query.Dir != nil {
			w.DB.Select("parent_id", "path").Find(&currentDirectory, "id = ?", *query.Dir)
		}
		if currentDirectory.ParentID != nil {
			parent = int(*currentDirectory.ParentID)
		} else {
			parent = -1
		}

		w.DB.Raw(`with recursive tree as (
						select id, parent_id, path, name
						from files
						where id = ?
						union all
						select f.id, f.parent_id, f.path, f.name
						from files as f
						join tree as parent on parent.parent_id = f.id
					)
					select *
					from tree;`, *query.Dir).Scan(&path)

		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}

	} else {
		current = -1
		parent = -1
	}

	return c.Render(200, "files", &FileData{
		PageTitle:          "Files - " + w.Settings.SiteName,
		IsAdmin:            account.IsAdmin,
		Files:              append(directories, files...),
		Path:               path,
		CurrentDirectoryID: current,
		ParentDirectoryID:  parent,
	})
}
