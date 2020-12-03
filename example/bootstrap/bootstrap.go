package bootstrap

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/kenretto/crane"
	"github.com/kenretto/crane/i18n"
	"github.com/kenretto/crane/response"
	"github.com/kenretto/crane/validator"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	validate = validator.New(validator.Translator{Translator: en.New()}, validator.Translator{Translator: zh.New()})
	pilot    crane.ICrane
	err      error
)

func init() {
	pilot, err = crane.NewCrane("application.yaml")
	if err != nil {
		panic(err)
	}
	response.SetTranslator(i18n.NewBundle(language.Chinese).LoadFiles("locales/", "yaml", yaml.Unmarshal))
	binding.Validator = validate
}

func Pilot() crane.ICrane {
	return pilot
}

func Validator() *validator.Validator {
	return validate
}
