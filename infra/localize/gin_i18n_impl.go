package localize

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

const (
	defaultRootPath = "./i18n"
)

var (
	defaultLanguage       = language.English
	defaultAcceptLanguage = []language.Tag{
		defaultLanguage,
		language.Chinese,
	}

	defaultBundleConfig = &BundleCfg{
		RootPath:        defaultRootPath,
		AcceptLanguage:  defaultAcceptLanguage,
		DefaultLanguage: defaultLanguage,
	}
)

// defaultLangResolver ...
func defaultLangResolver(context *gin.Context, defaultLang string) string {
	lang := context.Query("lang")
	if lang != "" {
		return lang
	}

	lang = context.GetHeader("Accept-Language")
	if lang != "" {
		return lang
	}

	return defaultLang
}

type GinI18n struct {
	bundleConfig *BundleCfg

	bundle          *i18n.Bundle
	currentContext  *gin.Context
	localizerByLang map[string]*i18n.Localizer
	defaultLanguage language.Tag
	langResolver    LangResolver
}

func (i *GinI18n) setBundleConfig(cfg *BundleCfg) {
	i.bundleConfig = cfg
}

// getMessage get localize message by lang and messageID
func (i *GinI18n) getMessage(messageID interface{}) (string, error) {
	lang := i.langResolver(i.currentContext, i.defaultLanguage.String())
	localizer := i.getLocalizerByLang(lang)

	localizeItem, err := i.getLocalizeItem(messageID)
	if err != nil {
		return fmt.Sprint(messageID), err
	}

	message, err := localizer.Localize(localizeItem)
	if err != nil {
		return fmt.Sprint(messageID), err
	}

	return message, nil
}

// mustGetMessage ...
func (i *GinI18n) mustGetMessage(param interface{}) string {
	message, _ := i.getMessage(param)
	return message
}

func (i *GinI18n) setCurrentContext(ctx context.Context) {
	i.currentContext = ctx.(*gin.Context)
}

func (i *GinI18n) setBundle(cfg *BundleCfg) {
	bundle := i18n.NewBundle(cfg.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	i.bundle = bundle
	i.defaultLanguage = cfg.DefaultLanguage

	i.loadMessageFiles(cfg)
	i.setLocalizerByLang(cfg.AcceptLanguage)
}

func (i *GinI18n) setLangResolver(handler LangResolver) {
	i.langResolver = handler
}

// loadMessageFiles load all file localize to bundle
func (i *GinI18n) loadMessageFiles(config *BundleCfg) {
	for _, lang := range config.AcceptLanguage {
		path := config.RootPath + "/" + lang.String() + ".yaml"
		i.bundle.MustLoadMessageFile(path)
	}
}

// setLocalizerByLang set localizer by language
func (i *GinI18n) setLocalizerByLang(acceptLanguage []language.Tag) {
	i.localizerByLang = map[string]*i18n.Localizer{}
	for _, lang := range acceptLanguage {
		langStr := lang.String()
		i.localizerByLang[langStr] = i.newLocalizer(langStr)
	}

	// set defaultLanguage if it isn't exist
	defaultLang := i.defaultLanguage.String()
	if _, hasDefaultLang := i.localizerByLang[defaultLang]; !hasDefaultLang {
		i.localizerByLang[defaultLang] = i.newLocalizer(defaultLang)
	}
}

// newLocalizer create a localizer by language
func (i *GinI18n) newLocalizer(lang string) *i18n.Localizer {
	langDefault := i.defaultLanguage.String()
	langs := []string{
		lang,
	}

	if lang != langDefault {
		langs = append(langs, langDefault)
	}

	localizer := i18n.NewLocalizer(
		i.bundle,
		langs...,
	)
	return localizer
}

// getLocalizerByLang get localizer by language
func (i *GinI18n) getLocalizerByLang(lang string) *i18n.Localizer {
	acceptLangs := ParseAcceptLanguage(lang)

	for _, al := range acceptLangs {
		localizer, hasValue := i.localizerByLang[al.Lang]
		if hasValue {
			return localizer
		}
	}

	return i.localizerByLang[i.defaultLanguage.String()]
}

func (i *GinI18n) getLocalizeItem(param interface{}) (*i18n.LocalizeConfig, error) {
	switch paramValue := param.(type) {
	case *i18n.LocalizeConfig:
		return paramValue, nil
	default:
		entry := &i18n.LocalizeConfig{
			MessageID: fmt.Sprintf("%s", param),
		}
		return entry, nil
	}
}
