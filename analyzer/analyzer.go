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
)

type Analyzer struct {
	context *context.Context
	master string
	downloadPath string
	stop bool
	maxProcess int
	interval int
	sembox chan string
}

type Analyzers []*Analyzer

func New(context *context.Context, master, downloadPath string, maxProcess int, interval int) (analyzer *Analyzer) {
	analyzer = &Analyzer{context:context, master:master, downloadPath:downloadPath}
	// 初始化并发数
	if maxProcess > 0 {
		analyzer.maxProcess = maxProcess
	} else {
		analyzer.maxProcess = 3
	}
	// 初始化每一次下载的时间间隔
	if interval > 0 {
		analyzer.interval = interval
	} else {
		analyzer.interval = 1
	}
	return
}

func (analyzers Analyzers) Init() (err error) {
	for _,analyzer := range analyzers {
		analyzer.sembox = make(chan string, analyzer.maxProcess)
		for i :=0; i < analyzer.maxProcess; i++ {
			analyzer.sembox <- msg.OK
		}
		analyzer.stop = true
		go analyzer.work()
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

	return conn,nil
}

/**
 * 与调度器建立链接
 */
func (analyzer *Analyzer) work() {
	for {
		if analyzer.stop == true {
			time.Sleep(1 * time.Second)
			continue
		}

		if msg.OK == <- analyzer.sembox {
			conn,err := connect(analyzer.master)
			if err != nil {
				log.Println("分析器链接调度器失败")
				log.Println(err)
				continue
			}
			
			url,err := getDownloadUrl(conn)
			if err != nil {
				log.Println("获取下载链接失败")
				log.Println(err)
				continue
			}
			go crawlUrl(analyzer, url, analyzer.sembox)
		}
		time.Sleep(time.Duration(analyzer.interval) * time.Second)
		
	}
}

/**
 * 抓取url
 */
func crawlUrl(analyzer *Analyzer, url string, sembox chan string) {
	statusCode,htmlPath,redirects,err := download.DownloadHTML(url, analyzer.downloadPath)

	if err != nil {
		log.Println(url + "下载失败")
	} else {
		go sendRedirectsToScheduler(analyzer.master, url, redirects)
	}

	if statusCode != 200 {
		log.Println(url+"返回码不是200")
		WriteURLInfo(analyzer, url, statusCode, nil)
	} else {
		log.Println("开始分析url："+url+" html文件在："+htmlPath)
		document := analyzer.analyse(htmlPath)
		WriteURLInfo(analyzer, url, statusCode, document)
	}

	sembox <- msg.OK
}

/**
 * 从调度器获取需要下载分析的url
 */
func getDownloadUrl(conn *net.TCPConn) (url string, err error) {
	defer conn.Close()
	_,err = socket.Write(conn, []byte("analyzer_ready"))
	if err != nil {
		return
	}
	result, err := socket.Read(conn)
	if err != nil {
		return
	}

	url = string(result)
	return
}

/**
 * 更新数据库中url和对应关键字的信息
 */
func WriteURLInfo(analyzer *Analyzer, url string, statusCode int, document *Document) {
	component,_ := analyzer.context.GetComponent("database")
	database := component.(*db.DatabaseConfig)
	sql,_ := database.Open()

	now := time.Now().Unix()
	var urlId int64
	err := sql.QueryRow("SELECT id FROM url_infos WHERE url=?", url).Scan(&urlId)
	if err == nil {
		// 更新url信息
		stmtIns, err := sql.Prepare("UPDATE url_infos SET statuscode=?, title=?, last_modified=?, description=? WHERE url=?")
		if err != nil {
			log.Println(err)
			// return
		}
		defer stmtIns.Close()
		if statusCode == 200 {
			stmtIns.Exec(statusCode, document.Title, now, document.MainContent, url)
		} else {
			stmtIns.Exec(statusCode, "", now, "", url)
			return
		}

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
			if i>= len(document.Keywords) {
				break
			}
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

		if i<len(ids) {
			for _,id := range ids {
				stmtIns, _ := sql.Prepare("DELETE FROM keywords WHERE id=?")
				stmtIns.Exec(id)
			}
		}
	} else {
		// 插入新url信息
		stmtIns, err := sql.Prepare("INSERT INTO url_infos(url, statuscode, title, last_modified, description) VALUES( ?, ?, ?, ?, ? )")
		if err != nil {
			log.Println(err)
			// return
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
	if len(redirects) ==0 {
		return
	}
	downloadResultMsg := msg.DownloadResultMsg{URL:url, Redirects:redirects}

	conn,err := connect(master)
	if err != nil {
		log.Println("发送redirect，下载器来链接调度器失败")
		log.Println(err)
		return
	}
	defer conn.Close()

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