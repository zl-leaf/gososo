package scheduler
import(
	"net"
	"log"
	"strings"
	"encoding/json"
	"time"

	"github.com/zl-leaf/gososo/scheduler/pool"
	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/socket"
	"github.com/zl-leaf/gososo/utils/db"
)

type Scheduler struct {
	context *context.Context
	port string
	listener *net.TCPListener
	stop bool
	analyzerPool *pool.AnalyzerPool
	maxTotal int64
}

func New(context *context.Context, port string, maxTotal int64) (scheduler *Scheduler){
	scheduler = &Scheduler{context:context, port:port}
	scheduler.analyzerPool = pool.NewAnalyzerPool()
	if maxTotal > 0 {
		scheduler.maxTotal = maxTotal
	} else {
		scheduler.maxTotal = -1
	}
	return
}

func (scheduler *Scheduler) Init() (err error){
	go scheduler.listenConnect()
	return
}

func (scheduler *Scheduler) Start() (err error){
	scheduler.initURLQueue()
	scheduler.stop = false
	go scheduler.dispatch()
	return
}

func (scheduler *Scheduler) Stop() (err error) {
	scheduler.stop = true
	return
}

/**
 * 初始化url抓取队列
 */
func (scheduler *Scheduler) initURLQueue() {
	component,_ := scheduler.context.GetComponent("database")
	database := component.(*db.DatabaseConfig)
	sql,_ := database.Open()

	rows, err := sql.Query("SELECT url FROM url_infos")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	analyseQueue.Clear()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url);err != nil {
			log.Println(err)
		}
		analyseQueue.Add(url)
	}
}

/**
 * 接收分析器的信息
 */
func (scheduler *Scheduler) listenConnect() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+scheduler.port)
	if err != nil {
		return
	}
	scheduler.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	listener := scheduler.listener
	for {
		time.Sleep(1 * time.Second)
		if scheduler.stop == true {
			continue
		}
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go scheduler.handle(conn)
	}
	return
}

func (scheduler *Scheduler) handle(conn net.Conn) {
	data,err := socket.Read(conn)
	if err != nil {
        return
    }

    message := strings.ToLower(string(data))
    switch {
    	case message==msg.ANALYZER_READY:
			_,err = socket.Write(conn, []byte(msg.OK))
			if err != nil {
				log.Println("发送接受准备确认信息出错")
				break
			}
			data,err := socket.Read(conn)
			if err != nil {
				log.Println("接收url匹配信息出错")
				break
			}
    		log.Printf("%s分析器准备就绪，负责分析%s\n", conn.RemoteAddr(), string(data))
    		scheduler.analyzerPool.Add(string(data), conn)

    	case message==msg.DOWNLOAD_OK:
    		_,err = socket.Write(conn, []byte(msg.OK))
    		if err != nil {
    			log.Println("发送接收新url信息时候出错")
    			break
    		}
    		data,err := socket.Read(conn)
    		if err != nil {
    			log.Println("接收redirect时候出错")
    			break
    		}

    		var downloadResultMsg msg.DownloadResultMsg
    		err = json.Unmarshal(data, &downloadResultMsg)
    		if err != nil {
    			log.Println("解析download result时候出错")
    			break
    		}

    		log.Println(downloadResultMsg.URL + "下载完成")
    		addRedirectURLs(downloadResultMsg.Redirects)
    }
}
