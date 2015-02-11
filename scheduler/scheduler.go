package scheduler
import(
	"net"
	"log"
	"./pool"
	"../utils/socket"
)

type Scheduler struct {
	port string
	listener *net.TCPListener
	stop bool
	downloaderPool *pool.DownloaderPool
}

func New(port string) (scheduler *Scheduler){
	scheduler = &Scheduler{port:port}
	scheduler.downloaderPool = pool.NewDownloaderPool()
	return
}

func (scheduler *Scheduler)Start() (err error){
	go scheduler.listen()
	go scheduler.dispatch()
	return
}

func (scheduler *Scheduler)Stop() {
	scheduler.listener.Close()
	scheduler.stop = true
	return
}

/**
 * 接收下载器和分析器的信息
 */
func (scheduler *Scheduler)listen() {
	log.Println("调度器开始监听")
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

func (scheduler *Scheduler)handle(conn net.Conn) {
	data,err := socket.Read(conn)
	if err != nil {
        return
    }
    if string(data) == "downloader_ready" {
    	log.Println("一个下载器准备就绪")
    	scheduler.downloaderPool.Add(conn)
    }
}
