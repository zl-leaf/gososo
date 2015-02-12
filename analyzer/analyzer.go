package analyzer
import(
	"net"
	"log"
	"time"
	"encoding/json"
	"../msg"
	"../utils/socket"
)

type Analyzer struct {
	master string
	dictionaryPath string
	stopwordsPath string
	stop bool
}

func New(master, dictionaryPath, stopwordsPath string) (analyzer *Analyzer) {
	analyzer = &Analyzer{master:master, dictionaryPath:dictionaryPath, stopwordsPath:stopwordsPath}
	return
}

func (analyzer *Analyzer)Start() (err error) {
	go analyzer.ready()
	return
}

func (analyzer *Analyzer)Stop() {
	analyzer.stop = true
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
func (analyzer *Analyzer)ready() {
	for {
		if analyzer.stop == true {
			break
		}
		time.Sleep(5 * time.Second)
		conn,err := connect(analyzer.master)
		if err != nil {
			log.Println("分析器链接调度器失败")
			log.Println(err)
			continue
		}
		_,err = socket.Write(conn, []byte("analyzer_ready"))
		if err != nil {
			log.Println("发送准备信息失败")
			log.Println(err)
			continue
		}
		data, err := socket.Read(conn)
		if err != nil {
			log.Println("下载器读取信息失败")
			log.Println(err)
			continue
		}

		var analyseOrderMsg msg.AnalyseOrderMsg
		err = json.Unmarshal(data, &analyseOrderMsg)
		if err != nil {
			log.Println("解析analyse命令时候出错")
			break
		}

		log.Println("开始分析url："+analyseOrderMsg.URL+" html文件在："+analyseOrderMsg.Path)
		analyzer.analyse(analyseOrderMsg.Path)
		
	}
}