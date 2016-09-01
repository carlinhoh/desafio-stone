package main

import (
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"sync"
)

type DatabaseFactory struct {
	DBConfig mysql.Config
	DB *sql.DB
}
var dbFactory *DatabaseFactory = new(DatabaseFactory)

var once sync.Once

func (dbFactory *DatabaseFactory) GetInstance() (*sql.DB, error){
	var err error
	once.Do(func() {
		if dbFactory.DB != nil{
			err = nil
		} else {
			dbFactory.DB, err = sql.Open("mysql", dbFactory.DBConfig.FormatDSN())
		}
	})
	return dbFactory.DB, err
}