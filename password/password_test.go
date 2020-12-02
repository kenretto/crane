package password

import (
	"github.com/kenretto/crane/configurator"
	"testing"
)

func TestPassword_Hash(t *testing.T) {
	var pwd = new(Password)
	var c, err = configurator.NewConfigurator("testdata/password.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add("password", pwd)
	if !pwd.Instance().Verify("hello", "world", "$2a$10$b2tatYGfdgjOqfFoVWNvWum47N45blcGn/HUrxc08oFtfkBTXQGGa") {
		t.Error("password error")
	}
}
