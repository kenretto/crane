package orm

import (
	"github.com/kenretto/crane/configurator"
	"testing"
	"time"
)

func TestLoader_DB(t *testing.T) {
	var loader = new(Loader)
	var c, err = configurator.NewConfigurator("testdata/database.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add("database", loader)

	var name string
	t.Log(loader.DB().Table("area").Where("area_id = ?", 1).Pluck("area_name", &name).RowsAffected)
	t.Log(name)
	time.Sleep(10 * time.Second)
	t.Log(loader.DB().Table("area").Where("area_id = ?", 1).Pluck("area_name", &name).RowsAffected)
	t.Log(name)
}
