package main
import(
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"strings"
	"log"
	"time"

	"github.com/zl-leaf/gososo/msg"
)

var searchUrl string = "http://localhost:9101/ajaxsearch/"

func main() {
	http.HandleFunc("/",indexHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
    	log.Fatal(err)
    }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	wd := r.Form.Get("wd")
	redirect := r.Form.Get("redirect")

	if strings.TrimSpace(redirect) != "" {
		// 跳转
		hit := r.Form.Get("hit") // 命中的关键字
		if strings.TrimSpace(hit) != "" {
			historyCookie,err := r.Cookie("hit")
			if err == nil && strings.TrimSpace(historyCookie.Value) != "" {
				history,err := url.QueryUnescape(historyCookie.Value)
				if err == nil {
					hitSlice := strings.Split(hit, ",")
					historySlice := strings.Split(history, ",")

					i := len(hitSlice)
					for j,his := range historySlice {
						if i+j > 100 {
							break
						}
						hit = hit + "," + his
					}
					
				}
			}
			hit = url.QueryEscape(hit)
			expire := time.Now().AddDate(0, 0, 30)
	    		cookie := http.Cookie{Name: "hit", Value: hit, Path: "/", Expires: expire, MaxAge: 30*86400}
	    		http.SetCookie(w, &cookie)
		}
		
    		http.Redirect(w, r, redirect, http.StatusFound)
		// w.Write([]byte(hit))
    		return
	}

	if strings.TrimSpace(wd) != "" {
		// 搜索
		gCookieJar, _ := cookiejar.New(nil)

		historyCookie,err := r.Cookie("hit")
		if err == nil {
			requestUrl,_ := url.Parse(searchUrl)
			cookies := []*http.Cookie{historyCookie}
			gCookieJar.SetCookies(requestUrl, cookies)
		}

		httpclient := http.Client{
		    	CheckRedirect: nil,
			Jar: gCookieJar,
		}

		resp, err := httpclient.Get(searchUrl+"?wd="+url.QueryEscape(wd))
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var searchResultMsg msg.SearchResultMsg
	    err = json.Unmarshal(body, &searchResultMsg)
	    if err != nil {
			w.Write([]byte("json error"))
			return
		}

		params := make(map[string]interface{})
		params["Data"] = searchResultMsg.Data

		t, err := template.ParseFiles("example/views/search.html")
		if err != nil {
	        log.Println(err)
	    }
		t.Execute(w, params)
	} else {
		// 首页
		t, err := template.ParseFiles("example/views/index.html")
		if err != nil {
	        log.Println(err)
	    }
	    t.Execute(w, nil)
	}
	
}

func redirectTo() {
	
}