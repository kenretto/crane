package util

import (
	"golang.org/x/text/language"
)

// GetLanguageLagByMobileCode 根据国家区号判断语言
func GetLanguageLagByMobileCode(code uint16) language.Tag {
	switch code {
	case 86, 886, 852, 853:
		return language.Chinese
	default:
		return language.English
	}
}
