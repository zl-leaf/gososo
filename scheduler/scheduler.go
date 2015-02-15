package scheduler
import(
	"net"
	"log"
	"strings"
	"encoding/json"
	"github.com/zl-leaf/gososo/scheduler/pool"
	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/socket"
)

type Scheduler struct {
	context *context.Context
	port string
	listener *net.TCPListener
	stop bool
	downloaderPool *pool.Pool
	analyzerPool *pool.Pool
}

func New(context *context.Context, port string) (scheduler *Scheduler){
	scheduler = &Scheduler{context:context, port:port}
	scheduler.downloaderPool = pool.NewDownloaderPool()
	scheduler.analyzerPool = pool.NewDownloaderPool()
	return
}

func (scheduler *Scheduler) Start() (err error){
	go scheduler.listen()
	go scheduler.dispatchDownload()
	go scheduler.dispatchAnalyse()
	return
}

func (scheduler *Scheduler) Stop() {
	scheduler.listener.Close()
	scheduler.stop = true
	return
}

/**
 * 接收下载器和分析器的信息
 */
func (scheduler *Scheduler) listen() {
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
		if scheduler.stop == true {
			break
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
    	case message==msg.DOWNLOAD_READY:
    		log.Printf("%s下载器准备就绪\n", conn.RemoteAddr())
    		scheduler.downloaderPool.Add(conn)

    	case message==msg.ANALYZER_READY:
    		log.Printf("%s分析器准备就绪\n", conn.RemoteAddr())
    		scheduler.analyzerPool.Add(conn)

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
    		addAnalyseURL(downloadResultMsg.URL, downloadResultMsg.Path)
    		addRedirectURLs(downloadResultMsg.Redirects)
    }
}
