package scheduler
import(
	"log"
	"../utils/socket"
	"../utils/queue"
)

var urlQueue *queue.Queue = queue.New()

func init() {
	urlQueue.Add("test1")
	urlQueue.Add("test2")
	urlQueue.Add("test3")
	urlQueue.Add("test4")
}

func (scheduler *Scheduler)dispatch() {
	for !urlQueue.Empty() {
		if scheduler.stop {
			break
		}

		e, err := urlQueue.Head()
		if err != nil {
			continue
		}
		u := e.Value.(string)
		conn := scheduler.downloaderPool.Get()
		log.Println("向下载器发送url："+u)
		_,err = socket.Write(conn, []byte(u))
		if err != nil {
			log.Println("发送url失败")
			log.Println(err)
		}
		conn.Close()
	}
}
