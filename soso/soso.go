package main
import(
	"log"
	"bufio"
	"os"
	"../sosoinit"
)

func main() {
	scheduler,downloaders := sosoinit.Sosoinit()

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