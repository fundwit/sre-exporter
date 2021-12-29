package localize

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func LocalizeMiddleware(i18nPath string) gin.HandlerFunc {
	return LocalizeMiddlewareWithCustomLangResolver(i18nPath, nil)
}

func LocalizeMiddlewareWithCustomLangResolver(i18nPath string, h LangResolver) gin.HandlerFunc {
	bundleCfg := &BundleCfg{
		RootPath:        i18nPath,
		AcceptLanguage:  []language.Tag{language.Chinese, language.English},
		DefaultLanguage: language.English,
	}

	options := []Option{}
	options = append(options, WithBundle(bundleCfg))

	if h != nil {
		options = append(options, WithCustomLangResolver(h))
	}

	return Localize(options...)
}
