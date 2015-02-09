package db

import (
	 _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

const(
	DRFAULT_HOST = "localhost"
	DEFAULT_PORT = "3306"
	DEFAULT_CHARSET = "utf8"
)

type DatabaseConfig struct {
	host string
	port string
	username string
	password string
	dbname string
	charset string
}

func New(config map[string]string) (dc *DatabaseConfig){
	host := config["host"]
	port := config["port"]
	username := config["username"]
	password := config["password"]
	dbname := config["dbname"]
	charset := config["charset"]

	if len(host) == 0 {
		host = DRFAULT_HOST
	}
	if len(port) == 0 {
		port = DEFAULT_PORT
	}
	if len(charset) == 0 {
		charset = DEFAULT_CHARSET
	}

	dc = &DatabaseConfig{host, port, username, password, dbname, charset}
	return
}

func (dc *DatabaseConfig)DSN() string{
	return dc.username+":"+dc.password+"@("+dc.host+":"+dc.port+")/"+dc.dbname+"?charset="+dc.charset
}

func (dc *DatabaseConfig)Open() (*sql.DB, error){
	db, err := sql.Open("mysql", dc.DSN())
	return db, err
}

func (dc *DatabaseConfig)CheckDBExist() (bool, error) {
	db,err := sql.Open("mysql", dc.DSN())
	if err != nil {
		return false,err
	}
	err = db.Ping()
	if err != nil {
		return false,err
	}
	return true,nil
}