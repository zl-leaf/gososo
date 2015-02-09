package main

import(
	"log"
	"../configure"
	"../utils/db"
)

const(
	SCHEDULER = "scheduler"
	DOWNLOADER = "downloader"
	DATABASE = "database"
)

func main() {
	config := configure.InitConfig("./config.ini")
	if dbConfig,exist := config.GetEntity(DATABASE);exist {
		checkDatabaseConfig(dbConfig)
	} else {
		log.Fatal("缺少数据库配置")
	}
	
}

func checkDatabaseConfig(es []*configure.Entity) {
	if len(es) > 1 {
		log.Fatal("数据库配置重复")
	}
	e := es[0]
	host := e.GetAttr("host")
	username := e.GetAttr("username")
	password := e.GetAttr("password")
	dbname := e.GetAttr("dbname")
	charset := e.GetAttr("charset")

	m := make(map[string]string)
	m["host"] = host
	m["username"] = username
	m["password"] = password
	m["dbname"] = dbname
	m["charset"] = charset

	databaseConfig := db.New(m)
	if exist,_ := databaseConfig.CheckDBExist();!exist {
		log.Fatal("数据库链接失败")
	}
}