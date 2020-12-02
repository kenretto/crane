package validator

import (
	"github.com/go-playground/locales"
	english "github.com/go-playground/locales/en"
	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/es"
	"github.com/go-playground/validator/v10/translations/fr"
	"github.com/go-playground/validator/v10/translations/id"
	"github.com/go-playground/validator/v10/translations/ja"
	"github.com/go-playground/validator/v10/translations/nl"
	"github.com/go-playground/validator/v10/translations/pt"
	"github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/go-playground/validator/v10/translations/ru"
	"github.com/go-playground/validator/v10/translations/tr"
	"github.com/go-playground/validator/v10/translations/zh"
)

// Translation translation interface
type Translation interface {
	locales.Translator
	RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error)
}

// Translator one-layer encapsulation of github.com/go-playground/locales
type Translator struct {
	locales.Translator
}

// ENTranslator English translator
func ENTranslator() Translator {
	return Translator{Translator: english.New()}
}

// ZHTranslator chinese translator
func ZHTranslator() Translator {
	return Translator{Translator: chinese.New()}
}

// RegisterDefaultTranslations locale
func (translator Translator) RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {
	switch translator.Locale() {
	case "en":
		return en.RegisterDefaultTranslations(v, trans)
	case "zh":
		return zh.RegisterDefaultTranslations(v, trans)
	case "es":
		return es.RegisterDefaultTranslations(v, trans)
	case "fr":
		return fr.RegisterDefaultTranslations(v, trans)
	case "id":
		return id.RegisterDefaultTranslations(v, trans)
	case "ja":
		return ja.RegisterDefaultTranslations(v, trans)
	case "nl":
		return nl.RegisterDefaultTranslations(v, trans)
	case "pt":
		return pt.RegisterDefaultTranslations(v, trans)
	case "pt_BR":
		return pt_BR.RegisterDefaultTranslations(v, trans)
	case "ru":
		return ru.RegisterDefaultTranslations(v, trans)
	case "tr":
		return tr.RegisterDefaultTranslations(v, trans)
	}
	return nil
}
