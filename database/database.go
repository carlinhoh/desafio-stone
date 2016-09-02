package database

import (
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"sync"
	"github.com/leoferlopes/desafio-stone/config"
)

type DatabaseFactory struct {
	DBConfig mysql.Config
	db *sql.DB
}

var dbFactory *DatabaseFactory

var once sync.Once

func GetInstance() (*sql.DB, error){
	var err error
	once.Do(func() {
		if dbFactory != nil && dbFactory.db != nil{
			err = nil
		} else {
			if dbFactory == nil{
				dbFactory = new(DatabaseFactory)

				var config *config.Config = &config.Settings

				dbFactory.DBConfig.Params = make(map[string]string)
				dbFactory.DBConfig.Params["charset"] = "utf8"
				dbFactory.DBConfig.Params["parseTime"] = "true"
				dbFactory.DBConfig.Addr = config.MySqlConfig.Address
				dbFactory.DBConfig.User = config.MySqlConfig.User
				dbFactory.DBConfig.Passwd = config.MySqlConfig.Password
				dbFactory.DBConfig.DBName = config.MySqlConfig.Schema
			}
			dbFactory.db, err = sql.Open("mysql", dbFactory.DBConfig.FormatDSN())
		}
	})
	return dbFactory.db, err
}