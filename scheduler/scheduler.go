package scheduler
import(
	"net"
	"log"
	"strings"
	"./pool"
	"../utils/socket"
)

type Scheduler struct {
	port string
	listener *net.TCPListener
	stop bool
	downloaderPool *pool.Pool
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

    msg := strings.ToLower(string(data))
    switch {
    	case msg=="downloader_ready":
    		log.Printf("%s下载器准备就绪\n", conn.RemoteAddr())
    		scheduler.downloaderPool.Add(conn)
    	case msg=="download_ok":
    		_,err = socket.Write(conn, []byte("ok"))
    		if err != nil {
    			log.Println("发送接收新url信息时候出错")
    			break
    		}
    		data,err := socket.Read(conn)
    		if err != nil {
    			log.Println("接受redirect时候出错")
    			break
    		}
    		redirects := strings.Split(string(data), "\n")
    		u := redirects[0]
    		log.Println(u + "下载完成")
    		redirects = redirects[1:]
    		addRedirectURLs(redirects)
    }
}
