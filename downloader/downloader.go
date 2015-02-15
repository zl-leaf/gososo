package downloader
import(
	"net"
	"log"
	"time"
	"encoding/json"

	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/socket"
)

type Downloader struct {
	context *context.Context
	master string
	downloadPath string
	stop bool
}

type Downloaders []*Downloader

func New(context *context.Context, master, downloadPath string) (downloader *Downloader){
	downloader = &Downloader{context:context, master:master, downloadPath:downloadPath}
	return
}

func (downloaders Downloaders) Init() (err error) {
	for _,downloader := range downloaders {
		go downloader.ready()
	}
	
	return
}

func (downloaders Downloaders) Start() (err error) {
	for _,downloader := range  downloaders {
		downloader.stop = false
	}
	
	return
}

func (downloaders Downloaders) Stop() (err error) {
	for _,downloader := range downloaders {
		downloader.stop = true
	}
	
	return
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
		statusCode,htmlPath,redirects,err := downloadHTML(url, downloader.downloadPath)
		if err != nil {
			log.Println(url + "下载失败")
		} else {
			downloadResultMsg := msg.DownloadResultMsg{URL:url, StatusCode:statusCode, Path:htmlPath, Redirects:redirects}
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