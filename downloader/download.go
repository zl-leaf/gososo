package downloader
import(
	"net/http"
	"net/url"
	"io/ioutil"
	"path/filepath"
	"os"
	"regexp"
)

func downloadHTML(u, downloadPath string) (statusCode int, htmlPath string, urls []string, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := string(b)

	htmlPath = downloadPath + resp.Request.URL.Host + resp.Request.URL.Path
	if htmlPath[len(htmlPath)-1:] == "/" {
		htmlPath = htmlPath + "index.html"
	}

	dir := filepath.Dir(htmlPath)
	err = os.MkdirAll(dir,0777)
	if err != nil {
		return
	}
	f,err := os.Create(htmlPath)
	if err != nil {
		return
	}
	_,err = f.WriteString(html)

	redirects := getRedirectURL(html)
	baseUrl := resp.Request.URL
	for _,redirect := range redirects {
		ref,err := url.Parse(redirect)
		if err == nil {
			u := baseUrl.ResolveReference(ref).String()
			urls = append(urls, u)
		}
	}
	return
}

func getRedirectURL(html string) (redirects []string) {
	redirects = make([]string, 0)

	// 去除注释
	notesRegexp := regexp.MustCompile(`(<\!\-\-)[\s\S]*?(\-\->)`)
	html = string(notesRegexp.ReplaceAll([]byte(html), []byte("")))

	hrefRegexp := regexp.MustCompile(`<a.*?href=\"(.*?[^\"])\".*?>.*?</a>`)
	match := hrefRegexp.FindAllStringSubmatch(html, -1)
	if match != nil {
 		for _,m := range match {
 			redirects = append(redirects, m[1])
 		}
	}
	return
}