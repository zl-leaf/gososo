package scheduler
import(
	"net"
	"net/url"
	"log"
	"time"

	"github.com/zl-leaf/gososo/utils/socket"
	"github.com/zl-leaf/gososo/utils/queue"
	"github.com/zl-leaf/gososo/scheduler/robots"

	"github.com/willf/bloom"
)

var analyseQueue *queue.Queue = queue.New()
var filter *bloom.BloomFilter = bloom.New(2700000, 5)
var mRobots *robots.Robots = robots.New("*")

func (scheduler *Scheduler) dispatch() {
	filter.ClearAll()
	finishPageNum := int64(0)
	for {
		time.Sleep(1 * time.Second)
		if scheduler.stop {
			break
		}


		if analyseQueue.Empty() {
			continue
		}
		e, err := analyseQueue.Head()
		if err != nil {
			continue
		}
		url := e.Value.(string)

		url = handleURL(url)
		if url == "" {
			continue
		}
		if filter.Test([]byte(url)) {
			// 可能已经抓取过
			continue
		}

		v,err := scheduler.analyzerPool.Get(url)
		if err != nil {
			continue
		}
		conn := v.(net.Conn)
		log.Println("向分析器发送url："+url)
		_,err = socket.Write(conn, []byte(url))
		if err != nil {
			log.Println("发送url失败")
			log.Println(err)
		}
		conn.Close()
		filter.Add([]byte(url))

		if scheduler.maxTotal > 0 {
			finishPageNum++
			if finishPageNum >= scheduler.maxTotal {
				log.Println("抓取达到上限")
				break
			}
		}
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
		u,err := url.Parse(redirect)
		if err == nil {
			robot := mRobots.GetRobot(u.Host)
			if robot.IsAllow(redirect) {
				analyseQueue.Add(redirect)
			}
		}
	}
}

/**
 * 处理URL，取出后面的Fragment
 */
func handleURL(e string) string {
	u,err := url.Parse(e)
	if err != nil {
		return ""
	}
	return u.RequestURI()
}
