package analyzer

import(
	"fmt"
	"io"
	"os"
	"log"
)

func (analyzer *Analyzer)analyse(f string) {
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	html := ""
	for {
		data := make([]byte, 500)
		count, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}
		html += string(data[:count])
	}

	document := &Document{}
	document.Init(analyzer.dictionaryPath, analyzer.stopwordsPath)
	document.LoadHTML(html)


	keywords := document.Keywords
	for _,kw := range keywords {
		fmt.Printf("%s %d %f\n", kw.Text, kw.Count, kw.Weight())
	}
}