package main
import(
	"log"
	"bufio"
	"os"

	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/sosoinit"
)

func main() {
	context := context.New()
	sosoinit.Sosoinit(context)

	scheduler,_ := context.GetService("scheduler")
	analyzers,_ := context.GetService("analyzers")

	// 启动监听
	if scheduler != nil {
		scheduler.Init()
	}
	if analyzers != nil {
		analyzers.Init()
	}

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
				if analyzers != nil {
					err := analyzers.Start()
					if err != nil {
						log.Println("分析器启动失败，错误如下")
						log.Println(err)
					}
				}
			case command=="http":
				if api,exist := context.GetService("api");exist {
					api.Init()
					api.Start()
				} else {
					log.Println("没有api配置")
				}
			case command=="stop":
				if scheduler != nil {
					scheduler.Stop()
				}
				if analyzers != nil {
					analyzers.Stop()
				}
			case command=="exit":
				os.Exit(1)
			default:
				log.Println("命令错误")
		}
	}
}