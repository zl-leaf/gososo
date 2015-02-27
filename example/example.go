package main
import(
	"net/http"
	"net/url"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"strings"
	"log"

	"github.com/zl-leaf/gososo/msg"
)

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

	if strings.TrimSpace(wd) != "" {
		resp, err := http.Get("http://localhost:9101/ajaxsearch/?wd="+url.QueryEscape(wd))
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
		t, err := template.ParseFiles("example/views/index.html")
		if err != nil {
	        log.Println(err)
	    }
	    t.Execute(w, nil)
	}
	
}