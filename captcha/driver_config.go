package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

// DriverConfig Configuration parameters of the captcha
type DriverConfig struct {
	CaptchaType string `mapstructure:"captcha_type"`
	Length      int    `mapstructure:"length"`
	Height      int    `mapstructure:"height"`
	Width       int    `mapstructure:"width"`
	NoiseCount  int    `mapstructure:"noise_count"`

	//string
	ShowLineOptions int         `mapstructure:"show_line_options"` // 2 or 4 or 8
	Source          string      `mapstructure:"source"`
	BackgroundColor *color.RGBA `mapstructure:"background"`
	Fonts           []string    `mapstructure:"fonts"`

	// audio
	Language string `mapstructure:"language"`

	// digit
	MaxSkew  float64 `mapstructure:"max_skew"`
	DotCount int     `mapstructure:"dot_count"`
}

// NewDriver generate a base64Captcha.Driver instance according to the configuration
func (config DriverConfig) NewDriver() base64Captcha.Driver {
	var driver base64Captcha.Driver
	switch config.CaptchaType {
	case "math":
		driver = base64Captcha.NewDriverMath(
			config.Height,
			config.Width,
			config.NoiseCount,
			config.ShowLineOptions,
			config.BackgroundColor,
			config.Fonts,
		)
	case "chinese":
		driver = base64Captcha.NewDriverChinese(
			config.Height,
			config.Width,
			config.NoiseCount,
			config.ShowLineOptions,
			config.Length,
			config.Source,
			config.BackgroundColor,
			config.Fonts,
		)
	case "digit":
		driver = base64Captcha.NewDriverDigit(config.Height, config.Width, config.Length, config.MaxSkew, config.DotCount)
	case "audio":
		driver = base64Captcha.NewDriverAudio(config.Length, config.Language)
	case "string":
		driver = base64Captcha.NewDriverString(
			config.Height,
			config.Width,
			config.NoiseCount,
			config.ShowLineOptions,
			config.Length,
			config.Source,
			config.BackgroundColor,
			config.Fonts,
		)
	}

	return driver
}
