package localize

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

type (
	// LangResolver ...
	LangResolver = func(context *gin.Context, defaultLang string) string

	// Option ...
	Option func(*GinI18n)
)

type BundleCfg struct {
	DefaultLanguage language.Tag
	AcceptLanguage  []language.Tag
	RootPath        string
}

// WithBundle ...
func WithBundle(config *BundleCfg) Option {
	return func(g *GinI18n) {
		g.setBundleConfig(config)
	}
}

// WithCustomLangResolver ...
func WithCustomLangResolver(f LangResolver) Option {
	return func(g *GinI18n) {
		g.setLangResolver(f)
	}
}

var atI18n *GinI18n

// newI18n ...
func newI18n(opts ...Option) {
	// init default value
	ins := &GinI18n{
		langResolver: defaultLangResolver,
		bundleConfig: defaultBundleConfig,
	}

	// overwrite default value by options
	for _, opt := range opts {
		opt(ins)
	}

	ins.setBundle(ins.bundleConfig)

	atI18n = ins
}

// Localize ...
func Localize(opts ...Option) gin.HandlerFunc {
	newI18n(opts...)
	return func(c *gin.Context) {
		atI18n.setCurrentContext(c)
	}
}

/*GetMessage get the i18n message
 param is one of these type: messageID, *i18n.LocalizeConfig
 Example:
	GetMessage("hello") // messageID is hello
	GetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": context.Param("name"),
			},
	})
*/
func GetMessage(param interface{}) (string, error) {
	return atI18n.getMessage(param)
}

/*MustGetMessage get the i18n message without error handling
  param is one of these type: messageID, *i18n.LocalizeConfig
  Example:
	MustGetMessage("hello") // messageID is hello
	MustGetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": context.Param("name"),
			},
	})
*/
func MustGetMessage(param interface{}) string {
	return atI18n.mustGetMessage(param)
}
