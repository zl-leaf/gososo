package main
import(
	"log"
	"bufio"
	"os"
	"../sosoinit"
)

func main() {
	context := sosoinit.Sosoinit()
	scheduler := context.Scheduler()
	downloaders := context.Downloaders()
	analyzers := context.Analyzers()
	log.Println("初始化完成")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch {
			case command=="start":
				if scheduler != nil {
					err := scheduler.Start()
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
				if scheduler != nil {
					scheduler.Stop()
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