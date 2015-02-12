package scheduler
import(
	"net"
	"log"
	"github.com/willf/bloom"
	"../utils/socket"
	"../utils/queue"
)

var urlQueue *queue.Queue = queue.New()
var filter *bloom.BloomFilter = bloom.New(2700000, 5)

func init() {
	urlQueue.Add("http://localhost/info/b.html")
	urlQueue.Add("http://localhost/info/a.html")
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
		url := e.Value.(string)
		if filter.Test([]byte(url)) {
			// 可能已经抓取过
			continue
		}

		conn := scheduler.downloaderPool.Get().(net.Conn)
		log.Println("向下载器发送url："+url)
		_,err = socket.Write(conn, []byte(url))
		if err != nil {
			log.Println("发送url失败")
			log.Println(err)
		}
		conn.Close()
		filter.Add([]byte(url))
	}
}

func addRedirectURLs(redirects []string) {
	for _,redirect := range redirects {
		if redirect == "" {
			continue
		}
		urlQueue.Add(redirect)
	}
}
