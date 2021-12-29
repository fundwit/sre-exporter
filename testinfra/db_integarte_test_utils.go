package testinfra

import (
	"context"
	"os"
	"sre-exporter/infra/persistence"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TestDatabase struct {
	TestDatabaseName string
	GormDB           *gorm.DB
}

// StartMysqlTestDatabase TEST_MYSQL_SERVICE=root:root@(127.0.0.1:3306)
func StartMysqlTestDatabase(baseName string) *TestDatabase {
	mysqlSvc := os.Getenv("TEST_MYSQL_SERVICE")
	if mysqlSvc == "" {
		mysqlSvc = "root:root@(127.0.0.1:3306)"
	}
	databaseName := baseName + "_test_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	dsn := "mysql://" + mysqlSvc + "/" + databaseName + "?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s"
	// create database (no conflict)
	if err := persistence.PrepareMysqlDatabase(dsn); err != nil {
		logrus.Fatalf("failed to prepare database %v\n", err)
	}

	db, err := persistence.StartGormDB(dsn)
	if err != nil {
		logrus.Fatalf("database connection failed %v\n", err)
	}

	return &TestDatabase{TestDatabaseName: databaseName, GormDB: db}
}

func StopMysqlTestDatabase(testDatabase *TestDatabase) {
	if testDatabase != nil || testDatabase.GormDB != nil {
		if err := testDatabase.GormDB.WithContext(context.Background()).Exec("DROP DATABASE " + testDatabase.TestDatabaseName).Error; err != nil {
			logrus.Println("failed to drop test database: " + testDatabase.TestDatabaseName)
		} else {
			logrus.Debugln("test database " + testDatabase.TestDatabaseName + " dropped")
		}

		// close connection
		persistence.StopGormDB(testDatabase.GormDB)
	}
}

func GormIntegrateTestSetup(t *testing.T, testDatabase **TestDatabase) {
	db := StartMysqlTestDatabase("skysight")
	*testDatabase = db
}

func GormIntegrateTestTeardown(t *testing.T, testDatabase *TestDatabase) {
	if testDatabase != nil {
		StopMysqlTestDatabase(testDatabase)
	}
}

// Example:
// 	func TestSomeIntegrate(t *testing.T) {
// 		gomega.RegisterTestingT(t)
// 		var testDatabase *testinfra.TestDatabase
//
// 		t.Run("gorm tracing should be ignored when parent span not found", func(t *testing.T) {
// 			defer gormIntegrateTestTeardown(t, testDatabase)
// 			gormIntegrateTestSetup(t, &testDatabase)
//
// 			Expect(testDatabase.GormDB.AutoMigrate(&TestResource{})).To(BeNil())
//
// 			// ...
// 		})
// 	}
