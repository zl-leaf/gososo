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
			case command=="exit":
				if scheduler != nil {
					scheduler.Stop()
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