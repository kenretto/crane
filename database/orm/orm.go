package orm

import (
	"github.com/kenretto/crane/unique"
	"gorm.io/gorm"
	"time"
)

type (
	// Database gorm supported data nesting, custom data table structure imported into the structure, will have these three fields by default
	Database struct {
		ID        uint64 `gorm:"primaryKey;column:id;" json:"id"`
		CreatedAt int64  `gorm:"column:created_at;index:created_at" json:"created_at"`
		UpdatedAt int64  `gorm:"column:updated_at;index:updated_at" json:"updated_at"`
	}
)

// BeforeCreate create hook, will force to set created_at and updated_at field, custom structure can re implement the method
func (db *Database) BeforeCreate(_ *gorm.DB) error {
	if db.ID == 0 {
		db.ID = unique.ID()
	}
	t := time.Now().Unix()
	db.CreatedAt, db.UpdatedAt = t, t
	return nil
}

// BeforeUpdate update hook
func (db *Database) BeforeUpdate(_ *gorm.DB) (err error) {
	t := time.Now().Unix()
	db.UpdatedAt = t
	return nil
}
