package downloader
import(
	"net/http"
	"net/url"
	"io/ioutil"
	"path/filepath"
	"os"
	"regexp"
)

func downloadHTML(u, downloadPath string) (urls []string, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := string(b)

	tempPath := downloadPath + resp.Request.URL.Host + resp.Request.URL.Path
	if tempPath[len(tempPath)-1:] == "/" {
		tempPath = tempPath + "index.html"
	}

	dir := filepath.Dir(tempPath)
	err = os.MkdirAll(dir,0777)
	if err != nil {
		return
	}
	f,err := os.Create(tempPath)
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
	var hrefRegexp = regexp.MustCompile(`<a.*?href=\"(.*?[^\"])\".*?>.*?</a>`)
	match := hrefRegexp.FindAllStringSubmatch(html, -1)
	if match != nil {
 		for _,m := range match {
 			redirects = append(redirects, m[1])
 		}
	}
	return
}