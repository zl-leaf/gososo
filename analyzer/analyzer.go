package analyzer
import(
	"net"
	"log"
	"time"
	"io"
	"os"
	"bufio"
	"strings"
	"encoding/json"

	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/socket"

	"github.com/huichen/sego"
)

type Analyzer struct {
	context *context.Context
	master string
	segmenter sego.Segmenter
	stopwords map[string]int
	stop bool
}

func New(context *context.Context, master, dictionaryPath, stopwordsPath string) (analyzer *Analyzer) {
	analyzer = &Analyzer{context:context, master:master}
	analyzer.segmenter.LoadDictionary(dictionaryPath)
	stopwords,err := getStopwrods(stopwordsPath)
	if err == nil {
		analyzer.stopwords = stopwords
	} else {
		analyzer.stopwords = make(map[string]int)
	}
	return
}

func getStopwrods(f string) (map[string]int,error){
	stopwords := make(map[string]int)
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		word, err := rd.ReadString('\n')
		word = strings.TrimSpace(word)
        if io.EOF == err {
            break
        }
		if err != nil {
            return nil, err
        }
        stopwords[word] = 0
	}
	return stopwords, nil
}

func (analyzer *Analyzer) Start() (err error) {
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
func (analyzer *Analyzer) ready() {
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