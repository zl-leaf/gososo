package analyzer

import(
	"io"
	"os"
	"log"
)

func (analyzer *Analyzer) analyse(f string) (document *Document) {
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

	document = &Document{}
	document.Init(analyzer.segmenter, analyzer.stopwords)
	document.LoadHTML(html)
	return
}