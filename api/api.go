package api
import(
	"net/http"
	"strings"
	"encoding/json"
	"log"

	"github.com/zl-leaf/gososo/context"
	"github.com/zl-leaf/gososo/msg"
	"github.com/zl-leaf/gososo/utils/db"
	"github.com/zl-leaf/gososo/utils/dictionary"

	"github.com/huichen/sego"
)

type Api struct {
	context *context.Context
	port string
	stop bool
}

func New(context *context.Context, port string) (api *Api){
	api = &Api{context:context, port:port}
	return
}

func (api Api) Init() (err error) {
	go api.listen()
	return
}

func (api Api) Start() (err error) {
	api.stop = false
	return
}

func (api Api) Stop() (err error) {
	api.stop = true
	return
}

/**
 * 接收信息
 */
func (api *Api) listen() {
	http.HandleFunc("/ajaxsearch/",api.ajaxSearchHandler)
    err := http.ListenAndServe(":"+api.port, nil)
    if err == nil {
    	log.Println("api服务器开启成功")
    } else {
    	log.Println("api服务器开启失败")
    }
	return
}

func (api *Api) ajaxSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	wd := r.Form.Get("wd")

	component,exist := api.context.GetComponent("dictionary")
	if !exist {
		w.Write([]byte("{result:-1, msg:'词典加载错误'}"))
		return
	}
	dictionary := component.(*dictionary.Dictionary)
	segmenter := dictionary.Sego()

	segments := segmenter.Segment([]byte(wd))
	words := sego.SegmentsToSlice(segments, true)
	searchResultMsg := searchFromDB(api, words)
	result,err := json.Marshal(searchResultMsg)
	if err != nil {
		w.Write([]byte("{result:-1, msg:'解析json出错'}"))
		return
	}
	w.Write([]byte(result))
	
}

func searchFromDB(api *Api, words []string) *msg.SearchResultMsg {
	searchResultMsg := &msg.SearchResultMsg{}

	if len(words) == 0 {
		searchResultMsg.Result = 1
		return searchResultMsg
	}

	wordsString := " ('" + strings.Join(words,"', '") + "') "

	component,exist := api.context.GetComponent("database")
	if exist {
		database := component.(*db.DatabaseConfig)
		sql,_ := database.Open()

		query := "select url,title,description,tmp.keywords from url_infos as url_info join "
		query += "(select url_id,group_concat(keyword) as keywords,sum(weight) as w from keywords where keyword in " + wordsString
		query += "group by url_id order by w desc)tmp "
		query += "on url_info.id=tmp.url_id"
		rows, err := sql.Query(query)

		if err != nil {
			searchResultMsg.Result = -1
			searchResultMsg.Msg = "服务端故障，查询失败"
			return searchResultMsg
		}
		defer rows.Close()

		searchResultMsg.Result = 1
		for rows.Next() {
			var url, title,description,keywords string
			if err := rows.Scan(&url, &title, &description, &keywords);err==nil {
				searchResultObj := &msg.SearchResultObj{URL:url, Title:title, Description:description, Keywords:keywords}
				searchResultMsg.Data = append(searchResultMsg.Data, searchResultObj)
			}
		}
	} else {
		searchResultMsg.Result = -1
		searchResultMsg.Msg = "服务端故障,数据库链接失败"
	}
	
	return searchResultMsg
}