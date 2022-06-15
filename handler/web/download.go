package web

import (
	"archive/zip"
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"github.com/myOmikron/hopfencloud/modules/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type DownloadArgs struct {
	ID uint `query:"id"`
}

func (w *Wrapper) Download(c echo.Context) error {
	var args DownloadArgs

	if err := c.Bind(&args); err != nil {
		//TODO: Render nicer template
		return c.String(400, "Bad request")
	}
	account, err := utils.GetAccount(c, w.DB)
	if err != nil {
		logger.Error(err.Error())
		return c.String(500, "Internal server error")
	}

	if args.ID == 0 {
		return c.String(400, "Bad ID")
	}

	// Check if user has access to id
	var count int64
	var file db.File
	w.DB.Find(&file, "id = ? AND owner_id = ? AND owner_type = 'accounts'", args.ID, account.ID).Count(&count)
	if count == 0 {
		return c.String(401, "You don't have access to this file")
	}

	if file.IsDirectory {
		tmpfile, err := ioutil.TempFile("", "")
		if err != nil {
			logger.Error(err.Error())
			return c.String(500, "Internal server error")
		}
		defer os.Remove(tmpfile.Name())

		zipWriter := zip.NewWriter(tmpfile)

		if err := filepath.Walk(file.Path, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				logger.Error(err.Error())
				return c.String(500, "Internal server error")
			}
			relPath := strings.TrimPrefix(filePath, filepath.Dir(file.Path))
			if zipFileWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
				Name:   relPath,
				Method: zip.Store,
			}); err != nil {
				logger.Error(err.Error())
				return c.String(500, "Internal server error")
			} else {
				if fsFile, err := os.Open(filePath); err != nil {
					logger.Error(err.Error())
					return c.String(500, "Internal server error")
				} else {
					if _, err = io.Copy(zipFileWriter, fsFile); err != nil {
						logger.Error(err.Error())
						return c.String(500, "Internal server error")
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}

		zipWriter.Close()
		return c.Attachment(tmpfile.Name(), file.Name+".zip")
	} else {
		return c.Attachment(file.Path, file.Name)
	}

}
