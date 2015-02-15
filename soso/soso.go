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

	schedulers,_ := context.GetService("schedulers")
	downloaders,_ := context.GetService("downloaders")
	analyzers,_ := context.GetService("analyzers")

	// 启动监听
	if schedulers != nil {
		schedulers.Init()
	}
	if downloaders != nil {
		downloaders.Init()
	}
	if analyzers != nil {
		analyzers.Init()
	}

	log.Println("初始化完成")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch {
			case command=="start":
				if schedulers != nil {
					err := schedulers.Start()
					if err != nil {
						log.Println("调度器启动失败，错误如下")
						log.Println(err)
					}
				}
				if downloaders != nil {
					err := downloaders.Start()
					if err != nil {
						log.Println("下载器启动失败，错误如下")
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
			case command=="exit":
				if schedulers != nil {
					schedulers.Stop()
				}
				if downloaders != nil {
					downloaders.Stop()
				}
				if analyzers != nil {
					analyzers.Stop()
				}
				os.Exit(1)
			default:
				log.Println("命令错误")
		}
	}
}