package persistence

import (
	"sre-exporter/infra/fail"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPrepareMysqlDatabase(t *testing.T) {
	RegisterTestingT(t)

	t.Run("only support mysql", func(t *testing.T) {
		Expect(PrepareMysqlDatabase("xxxx://aaa:bbb")).To(Equal(fail.ErrUnexpectedDatabase))
	})

	t.Run("error on invalid database url", func(t *testing.T) {
		Expect(PrepareMysqlDatabase("aa?bb/cc")).To(Equal(fail.ErrInvalidDatabaseUrl))
	})
}
