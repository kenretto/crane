package orm

import (
	"github.com/kenretto/crane/configurator"
	"github.com/sirupsen/logrus"
	"testing"
)

type area struct {
	ID   int    `gorm:"column:area_id;type:int"`
	Name string `gorm:"column:area_name;type:varchar(20)"`
}

func (area) TableName() string {
	return "area"
}

func TestLoader_DB(t *testing.T) {
	var loader = NewORM(logrus.NewEntry(logrus.New()))
	var c, err = configurator.NewConfigurator("testdata/database.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add(loader)

	var table area
	err = loader.DB().Migrator().AutoMigrate(&table)
	if err != nil {
		t.Error(err)
	}

	var name string
	t.Log(loader.DB().Table("area").Where("area_id = ?", 1).Pluck("area_name", &name).RowsAffected)
	t.Log(name)
}
