package orm

import (
	"github.com/kenretto/crane/unique"
	"gorm.io/gorm"
	"time"
)

type (
	// Database gorm 支持的数据嵌套, 自定义的数据表结构体导入该结构体，将默认拥有这三个字段
	Database struct {
		ID        uint64 `gorm:"primaryKey;column:id;" json:"id"`
		CreatedAt int64  `gorm:"column:created_at;index:created_at" json:"created_at"`
		UpdatedAt int64  `gorm:"column:updated_at;index:updated_at" json:"updated_at"`
	}
)

// BeforeCreate 创建数据前置操作, 将会强制设置 created_at 和 updated_at 字段, 自定义结构体可重新实现该方法
// 下级的结构体如果继承了 Database , 那么如果没有特别的需要不用重写 BeforeCreate 方法
func (db *Database) BeforeCreate(_ *gorm.DB) error {
	if db.ID == 0 {
		db.ID = unique.ID()
	}
	t := time.Now().Unix()
	db.CreatedAt, db.UpdatedAt = t, t
	return nil
}

// BeforeUpdate 更新数据前置操作, 自定义结构体可重新实现该方法
func (db *Database) BeforeUpdate(_ *gorm.DB) (err error) {
	t := time.Now().Unix()
	db.UpdatedAt = t
	return nil
}
