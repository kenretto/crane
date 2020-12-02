package password

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type ILogger interface {
	Error(args ...interface{})
}

type Password struct {
	Config struct {
		Token string `mapstructure:"token"`
		Cost  int    `mapstructure:"cost"`
	}
	l ILogger
}

func (pwd Password) OnChange(viper *viper.Viper) {
	_ = viper.Unmarshal(&pwd.Config)
}

func (pwd Password) Instance() Password {
	return pwd
}

// Hash 密码hash
func (pwd Password) Hash(token, password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s%s", pwd.Config.Token, password, token)), pwd.Config.Cost)
	if err != nil {
		pwd.l.Error(err)
	}

	return string(bytes)
}

// Verify 密码hash验证
func (pwd Password) Verify(token, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(fmt.Sprintf("%s%s%s", pwd.Config.Token, password, token)))
	return err == nil
}
