package localize

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseAcceptLanguage(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should ", func(t *testing.T) {
		Expect(ParseAcceptLanguage("zh-TW;q=0.4,zh-CN,zh;q=0.8,en-US;q=0.6")).
			To(Equal(AcceptLanguages{{"zh-CN", 1}, {"zh", 0.8}, {"en-US", 0.6}, {"zh-TW", 0.4}, {"en", 0}}))

		Expect(ParseAcceptLanguage("zh-TW;aa=0.4,zh-CN,en;q=xx")).
			To(Equal(AcceptLanguages{{"zh-CN", 1}, {"zh-TW", 0}, {"zh", 0}, {"en", 0}}))

		Expect(ParseAcceptLanguage("*")).
			To(Equal(AcceptLanguages{}))
	})
}
