package db

import (
	"time"

	"github.com/myOmikron/echotools/utilitymodels"
)

type FileVersion struct {
	utilitymodels.Common
	Path              string `gorm:"size:4096"`
	FileRelevantUntil time.Time
}

type File struct {
	utilitymodels.Common
	Name         string        `gorm:"size:256"`
	Path         string        `gorm:"size:4096"`
	FileVersions []FileVersion `gorm:"many2many:file__file_versions;"`
	IsDirectory  bool
	ParentID     *uint
	Parent       *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Permissions
	OwnerID        uint
	OwnerType      string
	InternalShares []InternalShare `gorm:"many2many:file__internal_shares;"`
	ExternalShares []ExternalShare `gorm:"many2many:file__external_shares;"`

	// Metadata
	Hash          string `gorm:"size:256"`
	Size          uint64
	FileCreatedAt time.Time
	FileUpdatedAt time.Time
}
