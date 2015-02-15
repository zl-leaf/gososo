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
	"github.com/zl-leaf/gososo/utils/db"

	"github.com/huichen/sego"
)

type Analyzer struct {
	context *context.Context
	master string
	segmenter sego.Segmenter
	stopwords map[string]int
	stop bool
}

type Analyzers []*Analyzer

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

func (analyzers Analyzers) Init() (err error) {
	for _,analyzer := range analyzers {
		go analyzer.ready()
	}
	
	return
}

func (analyzers Analyzers) Start() (err error) {
	for _,analyzer := range analyzers {
		analyzer.stop = false
	}
	
	return
}

func (analyzers Analyzers) Stop() (err error) {
	for _,analyzer := range analyzers {
		analyzer.stop = true
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
		document := analyzer.analyse(analyseOrderMsg.Path)
		analyzer.WriteURLInfo(analyseOrderMsg.URL, analyseOrderMsg.StatusCode, document)
		
	}
}

func (analyzer *Analyzer) WriteURLInfo(url string, statusCode int, document *Document) {
	component,_ := analyzer.context.GetComponent("database")
	database := component.(*db.DatabaseConfig)
	sql,_ := database.Open()

	var urlId int64
	err := sql.QueryRow("SELECT id FROM url_infos WHERE url=?", url).Scan(&urlId)
	if err != nil {
		log.Println(err)
	}

	now := time.Now().Unix()
	if urlId > 0 {
		// 更新url信息
		stmtIns, err := sql.Prepare("UPDATE url_infos SET statuscode=?, title=?, last_modified=?, description=? WHERE url=?")
		if err != nil {
			log.Println(err)
		}
		defer stmtIns.Close()
		stmtIns.Exec(statusCode, document.Title, now, document.MainContent, url)

		// 更新keywords
		rows,_ := sql.Query("SELECT id FROM keywords WHERE url_id=?", urlId)
		ids := make([]int64, 0)
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			ids = append(ids, id)
		}

		i := 0
		for i,id := range ids {
			keyword := document.Keywords[i]
			stmtIns, _ := sql.Prepare("UPDATE keywords SET keyword=?, weight=? WHERE id=?")
			stmtIns.Exec(keyword.Text, keyword.Weight(), id)
			i++
		}

		if i<len(document.Keywords) {
			stmtIns, _ = sql.Prepare("INSERT INTO keywords(keyword, weight, url_id) VALUES( ?, ?, ? )")
			for ;i<len(document.Keywords);i++ {
				keyword := document.Keywords[i]
				stmtIns.Exec(keyword.Text, keyword.Weight(), urlId)
			}
		}
		

	} else {
		// 插入新url信息
		stmtIns, err := sql.Prepare("INSERT INTO url_infos(url, statuscode, title, last_modified, description) VALUES( ?, ?, ?, ?, ? )")
		if err != nil {
			log.Println(err)
		}
		defer stmtIns.Close()

		res,_ := stmtIns.Exec(url, statusCode, document.Title, now, document.MainContent)

		// 插入keywords
		urlId,_ = res.LastInsertId()
		stmtIns, _ = sql.Prepare("INSERT INTO keywords(keyword, weight, url_id) VALUES( ?, ?, ? )")
		for _,keyword := range document.Keywords {
			stmtIns.Exec(keyword.Text, keyword.Weight(), urlId)
		}
	}
}