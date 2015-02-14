package main
import(
	"log"
	"bufio"
	"os"
	"../context"
	"../sosoinit"
	"../scheduler"
	"../downloader"
	"../analyzer"
)

func main() {
	context := context.New()
	sosoinit.Sosoinit(context)

	var schedul *scheduler.Scheduler
	var downloaders []*downloader.Downloader
	var analyzers []*analyzer.Analyzer

	if c,exist := context.GetComponent("scheduler");exist {
		schedul = c.(*scheduler.Scheduler)
	}
	if c,exist := context.GetComponent("downloaders");exist {
		downloaders = c.([]*downloader.Downloader)
	}
	if c,exist := context.GetComponent("analyzers");exist {
		analyzers = c.([]*analyzer.Analyzer)
	}

	log.Println("初始化完成")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch {
			case command=="start":
				if schedul != nil {
					err := schedul.Start()
					if err != nil {
						log.Println("调度器启动失败，错误如下")
						log.Println(err)
					}
				}
				if downloaders != nil {
					for _,d := range downloaders {
						err := d.Start()
						if err != nil {
							log.Println("下载器启动失败，错误如下")
							log.Println(err)
						}
					}
				}
				if analyzers != nil {
					for _, a := range analyzers {
						err := a.Start()
						if err != nil {
							log.Println("分析器启动失败，错误如下")
							log.Println(err)
						}
					}
				}
			case command=="exit":
				if schedul != nil {
					schedul.Stop()
				}
				if downloaders != nil {
					for _,d := range downloaders {
						d.Stop()
					}
				}
				os.Exit(1)
			default:
				log.Println("命令错误")
		}
	}
}