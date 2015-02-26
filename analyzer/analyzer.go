package analyzer
import(
	"net"
	"log"
	"time"
	"encoding/json"

	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/analyzer/download"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/socket"
	"github.com/zl-leaf/gososo/utils/db"

	"github.com/huichen/sego"
)

type Analyzer struct {
	context *context.Context
	master string
	downloadPath string
	segmenter sego.Segmenter
	stopwords map[string]int
	stop bool
}

type Analyzers []*Analyzer

func New(context *context.Context, master, downloadPath string) (analyzer *Analyzer) {
	analyzer = &Analyzer{context:context, master:master, downloadPath:downloadPath}
	
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
		result, err := socket.Read(conn)
		if err != nil {
			log.Println("下载器读取信息失败")
			log.Println(err)
			continue
		}

		url := string(result)

		statusCode,htmlPath,redirects,err := download.DownloadHTML(url, analyzer.downloadPath)

		if err != nil {
			log.Println(url + "下载失败")
		} else {
			go sendRedirectsToScheduler(analyzer.master, url, redirects)
		}

		log.Println("开始分析url："+url+" html文件在："+htmlPath)
		document := analyzer.analyse(htmlPath)
		analyzer.WriteURLInfo(url, statusCode, document)
		
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

func sendRedirectsToScheduler(master string, url string, redirects []string) {
	downloadResultMsg := msg.DownloadResultMsg{URL:url, Redirects:redirects}

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