package web

import (
	"fmt"
	"hash/adler32"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"github.com/myOmikron/hopfencloud/modules/utils"

	"github.com/labstack/echo/v4"
)

func (w *Wrapper) UploadPost(c echo.Context) error {
	now := time.Now()

	user, err := utils.GetAccount(c, w.DB)
	if err != nil {
		//TODO: Display error message
		return c.String(500, "Internal server error")
	}

	parent := c.FormValue("parent")
	file, err := c.FormFile("file")
	if err != nil {
		//TODO: Display error message
		return c.String(400, "Missing field file")
	}

	if file.Size == 0 {
		//TODO: Display error message
		return c.String(400, "Empty file")
	}

	// Check if parent is available for writing

	// Empty parent -> root
	var parentFile db.File
	var count int64
	if parent != "" {
		//TODO: Insert condition for group access
		w.DB.Find(
			&parentFile,
			"parent_id = ? AND is_directory = ? AND owner_id = ? AND owner_type = ?",
			parent, true, user.ID, "accounts",
		).Count(&count)

		// Invalid parent
		if count != 1 {
			//TODO: Display nicer error message
			return c.String(400, "Invalid parent directory")
		}
	}

	// Check if file already exists
	var parentFileID *uint
	if parent == "" {
		parentFileID = nil
	} else {
		parentFileID = &parentFile.ID
	}

	var existingFile db.File
	w.DB.Find(&existingFile, map[string]interface{}{
		"name":      file.Filename,
		"parent_id": parentFileID,
	}).Count(&count)
	if count != 0 {
		if existingFile.IsDirectory {
			//TODO: Display error message
			return c.String(500, "There's already a directory with that name")
		}
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create temp file
	temp, err := os.CreateTemp("", file.Filename+"*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	// Write to temp file
	if _, err = io.Copy(temp, src); err != nil {
		return err
	}

	// Calculate checksum
	checksum := adler32.New()
	if _, err := io.Copy(checksum, temp); err != nil {
		logger.Error(err.Error())
		//TODO: Display error message
		return c.String(500, "Internal server error")
	}
	if _, err := temp.Seek(0, 0); err != nil {
		logger.Error(err.Error())
		//TODO: Display error message
		return c.String(500, "Internal server error")
	}

	// Create db entry after new file is in place
	if count != 0 {
		versionPath := filepath.Join(
			utils.GetUserVersionsPath(user.ID, w.Config),
			fmt.Sprintf("%s.%d", file.Filename, now.Unix()),
		)

		// Move file from current/ to versions/
		if err := os.Rename(existingFile.Path, versionPath); err != nil {
			logger.Error(err.Error())
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}

		// Copy temp to current/
		newFile, err := os.Create(existingFile.Path)
		if err != nil {
			logger.Error(err.Error())
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}
		defer newFile.Close()

		if _, err := io.Copy(newFile, temp); err != nil {
			logger.Error(err.Error())
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}

		oldVersion := db.FileVersion{
			Path:              versionPath,
			FileRelevantUntil: now,
		}

		existingFile.Size = uint64(file.Size)
		existingFile.FileUpdatedAt = now
		existingFile.Hash = checksum.Sum32()

		w.DB.Save(&existingFile)
		if err := w.DB.Model(&existingFile).Association("FileVersions").Append(&oldVersion); err != nil {
			logger.Error(err.Error())
			return c.String(500, "Database error")
		}
	} else {
		// Copy temp to current/
		targetPath := utils.GetUserCurrentPath(user.ID, w.Config)

		if parentFileID == nil {
			targetPath = filepath.Join(targetPath, file.Filename)
		} else {
			targetPath = filepath.Join(targetPath, parentFile.Path, file.Filename)
		}

		newFile, err := os.Create(targetPath)
		if err != nil {
			logger.Error(err.Error())
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}
		defer newFile.Close()

		if _, err := io.Copy(newFile, temp); err != nil {
			logger.Error(err.Error())
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}

		w.DB.Create(&db.File{
			Name:          file.Filename,
			Path:          targetPath,
			Size:          uint64(file.Size),
			Hash:          checksum.Sum32(),
			FileCreatedAt: now,
			FileUpdatedAt: now,
			IsDirectory:   false,
			ParentID:      parentFileID,
			OwnerID:       user.ID,
			OwnerType:     "accounts",
		})
	}

	return c.Redirect(302, c.Request().Header.Get("Referer"))
}
