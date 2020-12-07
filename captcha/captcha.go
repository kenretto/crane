package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"sync"
)

type (
	// Loader Loader, this will be the first structure that should be derived from this package
	Loader struct {
		Driver DriverConfig `mapstructure:"driver"`
		Store  StoreConfig  `mapstructure:"store"`

		captcha *Captcha
		rw      sync.RWMutex
	}

	// Captcha Package github.com/mojocn/base64Captcha
	Captcha struct {
		*base64Captcha.Captcha
	}
)

func (loader *Loader) Node() string {
	return "captcha"
}

// NewCaptcha return Loader pointer
func NewCaptcha() *Loader {
	return new(Loader)
}

// OnChange Implemented the IConfig interface of the configurator package, finally register to configurator, will actively trigger this method when the configuration changes
func (loader *Loader) OnChange(viper *viper.Viper) {
	loader.rw.Lock()
	defer loader.rw.Unlock()
	_ = viper.Unmarshal(loader)
	loader.newInstance()
}

func (loader *Loader) newInstance() {
	loader.captcha = &Captcha{base64Captcha.NewCaptcha(loader.Driver.NewDriver(), loader.Store.NewStore())}
}

// Instance Get the verification code generation example, do not save this return value as a global variable
func (loader *Loader) Instance() *Captcha {
	loader.rw.RLock()
	defer loader.rw.RUnlock()
	return loader.captcha
}

// WithLogger Set up the logger required by this package
func (captcha *Captcha) WithLogger(logger ILogger) *Captcha {
	captcha.Store.(*RedisStore).SetLogger(logger)
	return captcha
}
