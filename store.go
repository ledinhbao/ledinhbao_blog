package main

import (
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"go.uber.org/zap"

	// log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func loadDatabase(dbconfig core.Config) (*gorm.DB, error) {
	var conn core.DatabaseConnection
	var err error

	// Any error here will lead to error on opening connection,
	// so just check it at one place.
	dialect, _ := dbconfig.StringValueForKey("dialect")
	databaseName, _ := dbconfig.StringValueForKey("database")
	username, _ := dbconfig.StringValueForKey("username")
	password, _ := dbconfig.StringValueForKey("password")
	host, _ := dbconfig.StringValueForKey("host")
	port, _ := dbconfig.StringValueForKey("port")

	conn, err = core.NewDatabaseConnection(dialect, databaseName, username, password, host, port)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialect, conn.ConnectionString())
	if err != nil {
		return nil, err
	}
	log.Info("database created",
		zap.String("dialect", dialect),
	)
	return db, nil
}
