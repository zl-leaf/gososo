package sosoinit
import(
	"../configure"
	"testing"
	"log"
)

func Test_Sosoinit(t *testing.T) {
	config := configure.InitConfig("../config.ini")

	if schedulerConfig,exist := config.GetEntity(SCHEDULER);exist {
		checkSchedulerConfig(schedulerConfig)
		scheduler := initScheduler(schedulerConfig)
		log.Println(scheduler)
	}

	if downloaderConfig,exist := config.GetEntity(DOWNLOADER);exist {
		checkDownloaderConfig(downloaderConfig)
		downloaders := initDownloaders(downloaderConfig)
		for _,downloader := range downloaders {
			log.Println(downloader)
		}
	}

	if dbConfig,exist := config.GetEntity(DATABASE);exist {
		checkDatabaseConfig(dbConfig)
	} else {
		log.Fatal("缺少数据库配置")
	}

	
	
}