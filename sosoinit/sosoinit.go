package sosoinit

import(
	"log"
	"../configure"
	"../scheduler"
	"../downloader"
	"../utils/db"
)

const(
	SCHEDULER = "scheduler"
	DOWNLOADER = "downloader"
	DATABASE = "database"

	MASTER = "master"
	PORT = "port"
	DOWNLOAD_PATH = "download_path"
)

func Sosoinit() (scheduler *scheduler.Scheduler, downloaders []*downloader.Downloader){
	config := configure.InitConfig("./config.ini")

	if schedulerConfig,exist := config.GetEntity(SCHEDULER);exist {
		checkSchedulerConfig(schedulerConfig)
		scheduler = initScheduler(schedulerConfig)
	}

	if downloaderConfig,exist := config.GetEntity(DOWNLOADER);exist {
		checkDownloaderConfig(downloaderConfig)
		downloaders = initDownloaders(downloaderConfig)
	}

	if dbConfig,exist := config.GetEntity(DATABASE);exist {
		checkDatabaseConfig(dbConfig)
	} else {
		log.Fatal("缺少数据库配置")
	}
	return
}

/**
 * 检查调度器的配置是否正确
 */
func checkSchedulerConfig(es []*configure.Entity) {
	if len(es) > 1 {
		//调度器只能有一个
		log.Fatal("scheduler配置重复")
	}
}

/**
 * 检查下载器的配置是否正确
 */
func checkDownloaderConfig(es []*configure.Entity) {
	for _,e := range es {
		if e.GetAttr(MASTER) == "" {
			log.Fatal("存在下载器没有对应master")
		}
	}
}

/**
 * 检查数据库配置是否正确
 */
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

/**
 * 初始化调度器
 */
func initScheduler(es []*configure.Entity) *scheduler.Scheduler{
	e := es[0]
	scheduler := scheduler.New(e.GetAttr(PORT))
	return scheduler
}

func initDownloaders(es []*configure.Entity) []*downloader.Downloader {
	downloaders := make([]*downloader.Downloader, 0)
	for _,e := range es {
		d := downloader.New(e.GetAttr(PORT), e.GetAttr(MASTER), e.GetAttr(DOWNLOAD_PATH))
		downloaders = append(downloaders, d)
	}
	return downloaders
}