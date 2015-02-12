package scheduler
import(
	"net"
	"log"
	"time"
	"encoding/json"
	"../msg"
	"github.com/willf/bloom"
	"../utils/socket"
	"../utils/queue"
)

var downloadQueue *queue.Queue = queue.New()
var analyseQueue *queue.Queue = queue.New()
var filter *bloom.BloomFilter = bloom.New(2700000, 5)

func init() {
	downloadQueue.Add("http://localhost/info/b.html")
	downloadQueue.Add("http://localhost/info/a.html")
}

func (scheduler *Scheduler)dispatchDownload() {
	for {
		if scheduler.stop {
			break
		}
		time.Sleep(1 * time.Second)

		if downloadQueue.Empty() {
			continue
		}

		e, err := downloadQueue.Head()
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

func (scheduler *Scheduler)dispatchAnalyse() {
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
		orderMsg := e.Value.(msg.AnalyseOrderMsg)

		conn := scheduler.analyzerPool.Get().(net.Conn)
		log.Println("向分析器发送url："+orderMsg.URL)
		result,err := json.Marshal(orderMsg)
		if err != nil {
			log.Println("analyse命令信息json化出错")
			continue
		}

		_,err = socket.Write(conn, result)
		if err != nil {
			log.Println("发送url失败")
			log.Println(err)
		}
		conn.Close()
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
		downloadQueue.Add(redirect)
	}
}

func addAnalyseURL(url,htmlPath string) {
	m := msg.AnalyseOrderMsg{URL:url, Path:htmlPath}
	log.Println("添加"+url+"到分析队列")
	analyseQueue.Add(m)
}