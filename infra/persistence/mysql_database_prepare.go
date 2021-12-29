package persistence

import (
	"os"
	"sre-exporter/infra/fail"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const mysqlDriverName = "mysql"

// PrepareMysqlDatabase  parameter example:
//   mysql://root:xxx@(test.xxx.com:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
//   or root:xxx@(test.xxx.com:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func PrepareMysqlDatabase(dsn string) error {
	driverName, driverArgs := splitName(dsn)
	if driverName != "" && driverName != mysqlDriverName {
		return fail.ErrUnexpectedDatabase
	}

	databaseName, rootDriverArgs, err := ExtractDatabaseName(driverArgs)
	if err != nil {
		return err
	}

	gormDB, err := gorm.Open(mysql.Open(rootDriverArgs), &gorm.Config{})
	if err != nil {
		return err
	}
	defer StopGormDB(gormDB)

	initSql := "CREATE DATABASE IF NOT EXISTS `" + databaseName + "` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;"

	if os.Getenv("GIN_MODE") == "debug" {
		gormDB = gormDB.Debug()
	}

	err = gormDB.Exec(initSql).Error
	if err != nil {
		return err
	}
	return nil
}
