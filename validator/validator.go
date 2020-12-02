package validator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"reflect"
	"strings"
	"sync"
)

var (
	_ binding.StructValidator = &Validator{} // 接口实现验证
)

// Validation custom validation interface
type Validation interface {
	Tag() string
	Validate(fl validator.FieldLevel) bool
	CallValidationEvenIfNull() bool
	Locale() string
	TranslateTmpl(ut ut.Translator) error
	TranslateParameters(ut ut.Translator, fe validator.FieldError) string
}

// Validator validator
type Validator struct {
	validate *validator.Validate
	config   map[string]map[string]string
	ut       *ut.UniversalTranslator
	mu       sync.RWMutex
}

// OnChange 配置发生变动时重新初始化
func (v *Validator) OnChange(viper *viper.Viper) {
	v.mu.Lock()
	defer v.mu.Unlock()
	_ = viper.Unmarshal(&v.config)
	v.init()
}

// ValidateStruct 验证结构体
//  Can be asserted as ValidationErrors,  if the err returned by validation is validator.ValidationErrors, Will be forced to ValidationErrors
func (v *Validator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := v.validate.Struct(obj); err != nil {
			if e, ok := err.(validator.ValidationErrors); ok {
				return ValidationErrors{ValidationErrors: e, validate: v}
			}
			return err
		}
	}

	return nil
}

// ValidationErrors 对验证错误的一层封装
type ValidationErrors struct {
	validator.ValidationErrors
	validate *Validator
}

// Translate localize errors after validation
func (ve ValidationErrors) Translate(locales ...string) validator.ValidationErrorsTranslations {
	var translator ut.Translator
	if len(locales) == 1 {
		translator = ve.validate.GetTranslator(locales[0])
	} else {
		translator = ve.validate.ut.GetFallback()
	}
	trans := make(validator.ValidationErrorsTranslations)
	for _, validationError := range ve.ValidationErrors {
		trans[validationError.Field()] = validationError.Translate(translator)
	}
	return trans
}

// Engine 获取验证器
func (v *Validator) Engine() interface{} {
	return v.validate
}

// New 获取验证器对象
// translators 必须传值, 第一个值将作为默认值
func New(translators ...Translation) *Validator {
	var v = new(Validator)
	var trans []locales.Translator
	for _, translator := range translators {
		trans = append(trans, translator)
	}

	v.ut = ut.New(translators[0], trans...)
	v.init()
	for _, translator := range translators {
		_ = v.RegisterDefaultTranslations(translator.Locale(), translator.RegisterDefaultTranslations)
	}
	return v
}

// RegisterDefaultTranslations custom translation
func (v *Validator) RegisterDefaultTranslations(locale string, translation func(v *validator.Validate, trans ut.Translator) (err error)) error {
	trans, _ := v.ut.GetTranslator(locale)
	err := translation(v.validate, trans)
	if err != nil {
		return err
	}
	return nil
}

// GetTranslator get translator
func (v *Validator) GetTranslator(locale string) ut.Translator {
	trans, _ := v.ut.GetTranslator(locale)
	return trans
}

// Bind gin controller
func (v *Validator) Bind(c *gin.Context, param interface{}) *ValidErrors {
	tags, _, err := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
	var locale string
	for _, tag := range tags {
		_, found := v.ut.GetTranslator(tag.String())
		if found {
			locale = tag.String()
			break
		}
	}

	if locale == "" {
		locale = v.ut.GetFallback().Locale()
	}

	if err != nil {
		var validErrors = newValidErrors()
		validErrors.Add("errors", "request error")
		return newValidErrors()
	}

	var validErrors = newValidErrors()
	if err := c.ShouldBind(param); err != nil {
		errs, ok := err.(ValidationErrors)
		if ok {
			for k, v := range errs.Translate(locale) {
				validErrors.Add(k, v)
			}
		} else {
			validErrors.Add("errors", err.Error())
		}
	}
	return validErrors
}

// RegisterTagNameFunc custom tag name
func (v *Validator) RegisterTagNameFunc(tags ...string) {
	v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		for _, tag := range tags {
			name := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}
		return field.Name
	})
}

// RegisterValidation set custom validation
func (v *Validator) RegisterValidation(validation Validation) error {
	trans, found := v.ut.GetTranslator(validation.Locale())
	if !found {
		return errors.New("translator not found")
	}
	_ = v.validate.RegisterValidation(validation.Tag(), validation.Validate, validation.CallValidationEvenIfNull())
	_ = v.validate.RegisterTranslation(validation.Tag(), trans, validation.TranslateTmpl, validation.TranslateParameters)
	return nil
}

// SetTagName set validator 获取验证规则的结构体 tag name
func (v *Validator) SetTagName(tag string) {
	v.validate.SetTagName(tag)
}

func (v *Validator) init() {
	v.validate = validator.New()
	v.RegisterTagNameFunc("form", "json", "xml")
}

// ValidErrors 验证之后的错误信息
type ValidErrors struct {
	ErrorsInfo map[string]string
	triggered  bool
}

// Add 添加错误信息
func (validErrors *ValidErrors) Add(key, value string) {
	validErrors.ErrorsInfo[key] = value
	validErrors.triggered = true
}

// IsValid 是否验证成功
func (validErrors *ValidErrors) IsValid() bool {
	return !validErrors.triggered
}

// Error return valid error's string
func (validErrors *ValidErrors) Error() string {
	return validErrors.String()
}

// String return valid error's string
func (validErrors *ValidErrors) String() (errString string) {
	for key, val := range validErrors.ErrorsInfo {
		errString += fmt.Sprintf("%s:%s\n", key, val)
	}
	return
}

func newValidErrors() *ValidErrors {
	return &ValidErrors{
		triggered:  false,
		ErrorsInfo: make(map[string]string),
	}
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
