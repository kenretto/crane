package password

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// ILogger logger
type ILogger interface {
	Error(args ...interface{})
}

// Password password
type Password struct {
	Config struct {
		Token string `mapstructure:"token"`
		Cost  int    `mapstructure:"cost"`
	}
	l ILogger
}

// OnChange ...
func (pwd Password) OnChange(viper *viper.Viper) {
	_ = viper.Unmarshal(&pwd.Config)
}

// Instance get password operation object
func (pwd Password) Instance() Password {
	return pwd
}

// Hash password hash
func (pwd Password) Hash(token, password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s%s", pwd.Config.Token, password, token)), pwd.Config.Cost)
	if err != nil {
		pwd.l.Error(err)
	}

	return string(bytes)
}

// Verify password hash verify
func (pwd Password) Verify(token, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(fmt.Sprintf("%s%s%s", pwd.Config.Token, password, token)))
	return err == nil
}
