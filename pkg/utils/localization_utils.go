package utils

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type LanguagesCode int

const (
	EN LanguagesCode = iota
	ZH               //Chinese
	ES               //Spanish
	ID               //Indonesian
	RU               //Russian
	UK               //English
	MY               // 缅甸语 Burmese
	TR               // 土耳其语 Turkish
	VI               // 越南语 Vietnamese
	BN               // 孟加拉语 Bengali
	HI               // 印地语 Hindi
	DE               //German
	UR               //Urdu
)

// Bengali, German, Spanish, Hindi, Indonesian, Burmese, Russian, Turkish, Urdu, Vietnamese, Chinese
var LanguagesCodeMap = map[string]LanguagesCode{
	"en": EN,
	"zh": ZH,
	"es": ES,
	"id": ID,
	"ru": RU,
	"uk": UK,
	"my": MY,
	"tr": TR,
	"vi": VI,
	"bn": BN,
	"hi": HI,
	"de": DE,
	"ur": UR,
}

var LocalizesCacheMap = map[string]*i18n.Localizer{}

// LocalizeStr localize string. if string localization not exist, key will be return
func LocalizeStr(key string, templateData map[string]interface{}, isPlural bool, languageCode string, localesBundle *i18n.Bundle) (string, error) {
	if _, ok := LanguagesCodeMap[languageCode]; !ok {
		fmt.Printf("localizeStr error, languageCode %s not support\n", languageCode)
		return key, nil
	}
	loc, ok := LocalizesCacheMap[languageCode]
	if !ok {
		localesBundle.MustLoadMessageFile(fmt.Sprintf("locales/active.%s.json", languageCode))
		loc = i18n.NewLocalizer(localesBundle, languageCode)
		LocalizesCacheMap[languageCode] = loc
	}
	count := 1
	if isPlural {
		count = 2
	}
	localizedStr, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templateData,
		PluralCount:  count,
	})
	if _, ok := err.(*i18n.MessageNotFoundErr); ok {
		fmt.Printf("LocalizeStr fail, key %s not existed\n", key)
		return key, nil
	}
	return localizedStr, err
}
