package testinfra

import (
	"database/sql/driver"
	"sre-exporter/infra/persistence"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyArgument) Match(v driver.Value) bool {
	return true
}

func SetUpMockSql() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mysqlConfig := mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}
	gormDB, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	persistence.ActiveGormDB = gormDB
	return gormDB, mock
}
