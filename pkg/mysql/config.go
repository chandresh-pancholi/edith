package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type DB struct {
	*sql.DB
	driver string
}

func NewDB() (*DB, error) {
	driver := viper.GetString("mysql.driver")
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))

	db, err := sql.Open(driver, dbSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.maxConnection"))
	db.SetMaxIdleConns(viper.GetInt("mysql.idleConnection"))

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db, driver}, nil
}
