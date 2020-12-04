package captcha

import (
	"context"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"github.com/kenretto/crane/configurator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"testing"
)

func TestLoader(t *testing.T) {
	var captcha = new(Loader)
	var c, err = configurator.NewConfigurator("testdata/captcha.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add("captcha", captcha)
	id, b64s, err := captcha.Instance().WithLogger(logrus.New()).Generate()
	if err != nil {
		t.Error(err)
	}
	data, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(b64s, "data:image/png;base64,", ""))
	if err != nil {
		t.Error(err)
	}
	_ = ioutil.WriteFile("testdata/captcha.png", data, 0666)
	v := captcha.Instance().Verify(id, captcha.Instance().WithLogger(logrus.New()).Store.Get(id, false), true)
	if !v {
		t.Error("captcha valid error")
	}
}

type MyLoader struct {
	Loader
	t *testing.T
}

func (my *MyLoader) OnChange(viper *viper.Viper) {
	my.Store.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		my.t.Log("hello, connected")
		return nil
	}
	my.Loader.OnChange(viper)
}

func TestImplLoader(t *testing.T) {
	var captcha = new(MyLoader)
	captcha.t = t
	var c, err = configurator.NewConfigurator("testdata/captcha.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add("captcha", captcha)
	id, b64s, err := captcha.Instance().WithLogger(logrus.New()).Generate()
	if err != nil {
		t.Error(err)
	}

	data, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(b64s, "data:image/png;base64,", ""))
	if err != nil {
		t.Error(err)
	}
	_ = ioutil.WriteFile("testdata/captcha.png", data, 0666)
	v := captcha.Instance().Verify(id, captcha.Instance().WithLogger(logrus.New()).Store.Get(id, false), true)
	if !v {
		t.Error("captcha valid error")
	}
}
