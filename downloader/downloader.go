package downloader
import(
	"net"
	"log"
	"time"
	"encoding/json"
	"../msg"
	"../utils/socket"
)

type Downloader struct {
	master string
	downloadPath string
	stop bool
}

func New(master, downloadPath string) (downloader *Downloader){
	downloader = &Downloader{master:master, downloadPath:downloadPath}
	return
}

func (downloader *Downloader) Start() (err error) {
	go downloader.ready()
	return
}

func (downloader *Downloader) Stop() {
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
func (downloader *Downloader) ready() {
	for {
		if downloader.stop == true {
			break
		}
		time.Sleep(5 * time.Second)
		conn,err := connect(downloader.master)
		if err != nil {
			log.Println("下载器链接调度器失败")
			log.Println(err)
			continue
		}
		_,err = socket.Write(conn, []byte(msg.DOWNLOAD_READY))
		if err != nil {
			log.Println("发送准备信息失败")
			log.Println(err)
			continue
		}
		result, err := socket.Read(conn)
		if err != nil {
			log.Println("下载器读取信息失败")
			log.Println(err)
			continue
		}

		url := string(result)
		htmlPath,redirects,err := downloadHTML(url, downloader.downloadPath)
		if err != nil {
			log.Println(url + "下载失败")
		} else {
			downloadResultMsg := msg.DownloadResultMsg{URL:url, Path:htmlPath, Redirects:redirects}
			sendRedirectsToScheduler(downloader.master, downloadResultMsg)
		}
	}
}

func sendRedirectsToScheduler(master string, downloadResultMsg msg.DownloadResultMsg) {
	conn,err := connect(master)
	if err != nil {
		log.Println("发送redirect，下载器来链接调度器失败")
		log.Println(err)
		return
	}

	_,err = socket.Write(conn, []byte(msg.DOWNLOAD_OK))
	if err != nil {
		log.Println("发送下载完成信息时候出错")
		return
	}
	data,err := socket.Read(conn)
	if err != nil || string(data) != msg.OK {
		log.Println("接收master确认信息出错")
		return
	}

	result,err := json.Marshal(downloadResultMsg)
	if err != nil {
		log.Println("download返回信息json化出错")
		return
	}
	_,err = socket.Write(conn, result)
	if err != nil {
		log.Println("发送redirect失败")
	}
}