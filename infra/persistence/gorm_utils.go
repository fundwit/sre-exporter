package persistence

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

var ActiveGormDB *gorm.DB

func StartGormDB(dsn string) (*gorm.DB, error) {
	_, args := splitName(dsn)
	gormDB, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	gormDB.Use(gormopentracing.New(gormopentracing.WithLogResult(false)))

	if os.Getenv("GIN_MODE") == "debug" {
		gormDB = gormDB.Debug()
	}

	return gormDB, nil
}

func StopGormDB(gormDB *gorm.DB) {
	if gormDB != nil {
		if conn, ok := gormDB.ConnPool.(interface{ Close() error }); ok {
			if err := conn.Close(); err != nil {
				logrus.Warnln("failed to close DB:", err)
			}
		}
	}
}
