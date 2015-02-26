package sosoinit

import(
	"log"
	"os"
	
	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/configure"
	"github.com/zl-leaf/gososo/scheduler"
	"github.com/zl-leaf/gososo/analyzer"
	"github.com/zl-leaf/gososo/utils/db"
	"github.com/zl-leaf/gososo/utils/dictionary"
	"github.com/zl-leaf/gososo/api"
)

const(
	SCHEDULER = "scheduler"
	DOWNLOADER = "downloader"
	ANALYZER = "analyzer"
	DATABASE = "database"
	API = "api"
	DICTIONARY = "dictionary"

	MASTER = "master"
	PORT = "port"
	DOWNLOAD_PATH = "download_path"
	DICTIONARY_PATH = "dictionary_path"
	STOPWORDS_PATH = "stopwords_path"
)

func Sosoinit(context *context.Context) {
	var scheduler *scheduler.Scheduler
	var analyzers analyzer.Analyzers
	var database *db.DatabaseConfig
	var api *api.Api
	var dictionary *dictionary.Dictionary

	config := configure.InitConfig("./config.ini")

	if schedulerConfig,exist := config.GetEntity(SCHEDULER);exist {
		checkSchedulerConfig(schedulerConfig)
		scheduler = initScheduler(context, schedulerConfig)
	}

	if analyzerConfig,exist := config.GetEntity(ANALYZER);exist {
		checkAnalyzerConfig(analyzerConfig)
		analyzers = initAnalyzers(context, analyzerConfig)
	}

	if dictionaryConfig,exist := config.GetEntity(DICTIONARY);exist {
		checkDictConfig(dictionaryConfig)
		dictionary = initDict(dictionaryConfig)
	}

	if dbConfig,exist := config.GetEntity(DATABASE);exist {
		checkDatabaseConfig(dbConfig)
		database = initDB(dbConfig)
	} else {
		log.Fatal("缺少数据库配置")
	}

	if apiConfig,exist := config.GetEntity(API);exist {
		checkApiConfig(apiConfig)
		api = initApi(context, apiConfig)
	}

	context.AddService("scheduler", scheduler)
	context.AddService("analyzers", analyzers)
	context.AddService("api", api)

	context.AddComponent("database", database)
	context.AddComponent("dictionary", dictionary)
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
 * 检查分析器的配置是否正确
 */
func checkAnalyzerConfig(es []*configure.Entity) {
	for i,e := range es {
		if e.GetAttr(MASTER) == "" {
			log.Fatal("存在分析器没有对应master")
		}
		if e.GetAttr(DOWNLOAD_PATH) != "" {
			dir := e.GetAttr(DOWNLOAD_PATH)

			fi, err := os.Stat(dir)
			if err != nil && !os.IsExist(err) || !fi.IsDir() {
				err := os.MkdirAll(dir,0777)
				if err != nil {
					log.Fatalf("第%d个分析器的下载路径无法生成\n", i)
				}
			}
		} else {
			log.Fatalf("第%d个分析器没有配置下载路径\n", i)
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
	if e == nil {
		return
	}
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

func checkDictConfig(es []*configure.Entity) {
	if len(es) > 1 {
		log.Fatal("dictionary配置重复")
	}

	e := es[0]
	if e == nil {
		return
	}

	if e.GetAttr(DICTIONARY_PATH) != "" {
		dictionary := e.GetAttr(DICTIONARY_PATH)
		fi, err := os.Stat(dictionary)
		if err != nil && !os.IsExist(err) || fi.IsDir() {
			log.Println("dictionary的词典路径错误")
		}
	}

	if e.GetAttr(STOPWORDS_PATH) != "" {
		stopwords := e.GetAttr(STOPWORDS_PATH)
		fi, err := os.Stat(stopwords)
		if err != nil && !os.IsExist(err) || fi.IsDir() {
			log.Println("dictionary的停用词典路径错误")
		}
	}
}

/**
 * 检查api的配置
 */
func checkApiConfig(es []*configure.Entity) {
	if len(es) > 1 {
		log.Fatal("api配置重复")
	}
}

/**
 * 初始化调度器
 */
func initScheduler(context *context.Context, es []*configure.Entity) *scheduler.Scheduler {
	if len(es) > 0 {
		e := es[0]
		scheduler := scheduler.New(context, e.GetAttr(PORT))
		return scheduler
	} else {
		return nil
	}
	
}

func initAnalyzers(context *context.Context, es []*configure.Entity) analyzer.Analyzers {
	analyzers := make(analyzer.Analyzers, 0)
	for _,e := range es {
		a := analyzer.New(context, e.GetAttr(MASTER), e.GetAttr(DOWNLOAD_PATH))
		analyzers = append(analyzers, a)
	}
	return analyzers
}

func initDB(es []*configure.Entity) *db.DatabaseConfig {
	if len(es) > 1 {
		log.Fatal("数据库配置重复")
	}
	e := es[0]
	if e == nil {
		return nil
	}
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
	return databaseConfig
}

func initDict(es []*configure.Entity) *dictionary.Dictionary {
	if len(es) > 1 {
		log.Fatal("dictionary配置重复")
	}
	e := es[0]
	if e == nil {
		return nil
	}

	dict := dictionary.New(e.GetAttr(DICTIONARY_PATH), e.GetAttr(STOPWORDS_PATH))
	return dict
}

func initApi(context *context.Context, es []*configure.Entity) *api.Api {
	if len(es) > 0 {
		e := es[0]
		a := api.New(context, e.GetAttr(PORT))
		return a
	} else {
		return nil
	}
}