package persistence

import (
	"errors"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDatabaseConfig(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return a DatabaseConfig instance when database url env is valid", func(t *testing.T) {
		os.Setenv(EnvDatabaseURL, "mysql://user.pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true")
		dsn, err := ParseDatabaseConfigFromEnv()

		Expect(err).To(BeZero())
		Expect(dsn).To(Equal("mysql://user.pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true"))
	})

	t.Run("should return err when database url env is not valid", func(t *testing.T) {
		os.Setenv(EnvDatabaseURL, "user.pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true")
		config, err := ParseDatabaseConfigFromEnv()

		Expect(err).To(Equal(errors.New(EnvDatabaseURL + " is not valid, a correct example like 'mysql://user:pwd@tcp(host:3306)/dbname?para=value'")))
		Expect(config).To(BeZero())
		Expect(strings.Contains(err.Error(), "is not valid, a correct example like")).To(BeTrue())
	})
}

func TestExtractDatabaseName(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should work correctly", func(t *testing.T) {
		var name, rootUrl string
		var err error
		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)/dbname?charset=utf8mb4")
		Expect(err).To(BeNil())
		Expect(name).To(Equal("dbname"))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4"))

		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4")
		Expect(err).To(BeNil())
		Expect(name).To(Equal(""))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4"))

		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)?charset=utf8mb4")
		Expect(err).To(BeNil())
		Expect(name).To(Equal(""))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)?charset=utf8mb4"))

		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)/dbname")
		Expect(err).To(BeNil())
		Expect(name).To(Equal("dbname"))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)/"))

		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)/")
		Expect(err).To(BeNil())
		Expect(name).To(Equal(""))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)/"))

		name, rootUrl, err = ExtractDatabaseName("root:P@4word@(test.xxxxx.com:3308)")
		Expect(err).To(BeNil())
		Expect(name).To(Equal(""))
		Expect(rootUrl).To(Equal("root:P@4word@(test.xxxxx.com:3308)"))

		// ...?.../...
		name, rootUrl, err = ExtractDatabaseName("root?abc/def")
		Expect(err).ToNot(BeNil())
		Expect(name).To(BeZero())
		Expect(rootUrl).To(BeZero())
	})
}

func TestSplitName(t *testing.T) {
	RegisterTestingT(t)
	t.Run("works as expected", func(t *testing.T) {
		d, n := splitName("aa://bbb/cc")
		Expect(d).To(Equal("aa"))
		Expect(n).To(Equal("bbb/cc"))

		d, n = splitName("aa:/bbb/cc")
		Expect(d).To(Equal(""))
		Expect(n).To(Equal("aa:/bbb/cc"))
	})
}
