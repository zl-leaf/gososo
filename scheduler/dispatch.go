package scheduler
import(
	"net"
	"log"
	"time"
	
	"github.com/zl-leaf/gososo/utils/socket"
	"github.com/zl-leaf/gososo/utils/queue"

	"github.com/willf/bloom"
)

var analyseQueue *queue.Queue = queue.New()
var filter *bloom.BloomFilter = bloom.New(2700000, 5)


func (scheduler *Scheduler) dispatch() {
	for {
		if scheduler.stop {
			break
		}
		time.Sleep(1 * time.Second)

		if analyseQueue.Empty() {
			continue
		}

		e, err := analyseQueue.Head()
		if err != nil {
			continue
		}
		url := e.Value.(string)
		if filter.Test([]byte(url)) {
			// 可能已经抓取过
			continue
		}

		conn := scheduler.analyzerPool.Get().(net.Conn)
		log.Println("向分析器发送url："+url)
		_,err = socket.Write(conn, []byte(url))
		if err != nil {
			log.Println("发送url失败")
			log.Println(err)
		}
		conn.Close()
		filter.Add([]byte(url))
	}
}

/**
 * 添加新url到待下载的url的队列
 */
func addRedirectURLs(redirects []string) {
	for _,redirect := range redirects {
		if redirect == "" {
			continue
		}
		log.Println("添加"+redirect+"到下载队列")
		analyseQueue.Add(redirect)
	}
}