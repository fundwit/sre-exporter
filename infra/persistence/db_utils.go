package persistence

import (
	"errors"
	"os"
	"sre-exporter/infra/fail"
	"strings"
)

const EnvDatabaseURL = "DATABASE_URL"

func ParseDatabaseConfigFromEnv() (string, error) {
	databaseUrl := os.ExpandEnv(os.Getenv(EnvDatabaseURL))

	if !strings.Contains(databaseUrl, "://") {
		return "", errors.New(EnvDatabaseURL + " is not valid, a correct example like 'mysql://user:pwd@tcp(host:3306)/dbname?para=value'")
	}

	return databaseUrl, nil
}

func ExtractDatabaseName(mysqlDriverArgs string) (string, string, error) {
	nameIndex := strings.IndexRune(mysqlDriverArgs, '/')
	paramsIndex := strings.IndexRune(mysqlDriverArgs, '?')

	// .../..?..
	if nameIndex > 0 && paramsIndex > nameIndex {
		return mysqlDriverArgs[nameIndex+1 : paramsIndex], mysqlDriverArgs[0:nameIndex+1] + mysqlDriverArgs[paramsIndex:], nil
	}
	// without /
	if nameIndex < 0 {
		return "", mysqlDriverArgs, nil
	}
	// with / and without ?
	if nameIndex > 0 && paramsIndex < 0 {
		return mysqlDriverArgs[nameIndex+1:], mysqlDriverArgs[0 : nameIndex+1], nil
	}

	// ..?../..
	return "", "", fail.ErrInvalidDatabaseUrl
}

func splitName(dsn string) (string, string) {
	var driverName string
	var driverArgs string
	if strings.Contains(dsn, "://") {
		parts := strings.Split(dsn, "://")
		driverName = strings.ToLower(parts[0])
		driverArgs = parts[1]
	} else {
		driverName = ""
		driverArgs = dsn
	}

	return driverName, driverArgs
}
