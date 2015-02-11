package downloader
import(
	"net"
	"log"
	"fmt"
	"time"
	"../utils/socket"
)

type Downloader struct {
	port string
	master string
	downloadPath string
	stop bool
}

func New(port, master, downloadPath string) (downloader *Downloader){
	downloader = &Downloader{port:port, master:master, downloadPath:downloadPath}
	return
}

func (downloader *Downloader)Start() (err error) {
	go downloader.ready()
	return
}

func (downloader *Downloader)Stop() {
	downloader.stop = true
}

/**
 * 链接调度器
 */
func connect(ip string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip)
	if err != nil {
		return nil,err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil,err
	}
	conn.SetKeepAlive(true)

	return conn,nil
}


/**
 * 发送准备就绪信息到调度器
 */
func (downloader *Downloader)ready() {
	log.Println("下载器开始链接调度器")
	for {
		if downloader.stop == true {
			break
		}
		time.Sleep(5 * time.Second)
		conn,err := connect(downloader.master)
		if err != nil {
			log.Println("下载器来链接调度器失败")
			log.Println(err)
		}
		_,err = socket.Write(conn, []byte("downloader_ready"))
		if err != nil {
			log.Println("发送准备信息失败")
			log.Println(err)
		}
		result, err := socket.Read(conn)
		if err != nil {
			log.Println("下载器读取信息失败")
			log.Println(err)
		}
		fmt.Println("下载器接收到调度器的信息："+string(result))
	}
}