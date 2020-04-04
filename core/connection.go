package core

import (
	"errors"
	"fmt"
)

// DatabaseConnectionError presents error when creating DatabaseConnection object
type DatabaseConnectionError struct {
	Message string
}

// DatabaseConnection provide ConnectionString() string, method to open database
type DatabaseConnection interface {
	ConnectionString() string
}

type sqlite3Connection struct {
	filepath string
}

type mysqlConnection struct {
	username     string
	password     string
	databaseName string
}
type postgresConnection struct {
	name, user, pass, host, port string
}

func (conn sqlite3Connection) ConnectionString() string {
	return conn.filepath
}
func (conn mysqlConnection) ConnectionString() string {
	return fmt.Sprintf("%s:%s@/%s?parseTime=true", conn.username, conn.password, conn.databaseName)
}
func (conn postgresConnection) ConnectionString() string {
	return fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s",
		conn.name, conn.user, conn.pass,
		conn.host, conn.port,
	)
}

func (err DatabaseConnectionError) Error() string {
	return err.Message
}

// NewDatabaseConnection create database connection object, based on dialect.
// 	- "sqlite3"  requires 1 arg : databaseName
// 	- "mysql"    requires 3 args: databaseName, username, password
// 	- "postgres" requires 5 args: databaseName, username, password, host, port
func NewDatabaseConnection(dialect string, args ...string) (DatabaseConnection, error) {
	switch dialect {
	case "sqlite", "sqlite3":
		if len(args) < 1 {
			return nil, DatabaseConnectionError{Message: "Missing filepath for sqlite connection"}
		}
		return &sqlite3Connection{
			filepath: args[0],
		}, nil
	case "mysql":
		if len(args) < 3 {
			return nil, DatabaseConnectionError{Message: "mysql connection requires username, password and databaseName"}
		}
		return &mysqlConnection{
			username:     args[1],
			password:     args[2],
			databaseName: args[0],
		}, nil
	case "postgres":
		if len(args) < 5 {
			return nil, DatabaseConnectionError{Message: "postgres connection requires databaseName, username, password, host, port"}
		}
		return &postgresConnection{
			name: args[0], user: args[1], pass: args[2],
			host: args[3], port: args[4],
		}, nil
	}
	if dialect == "" {
		return nil, errors.New("NewDatabaseConnection failed: empty dialect")
	}
	return nil, fmt.Errorf("NewDatabaseConnection failed: unsupported dialect %s", dialect)
}
