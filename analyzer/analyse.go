package analyzer

import(
	"io"
	"os"
	"log"

	"github.com/zl-leaf/gososo/utils/dictionary"
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
	component,exist := analyzer.context.GetComponent("dictionary")
	if component == nil {
		log.Println("dictionary error")
	}
	if exist {
		dictionary := component.(*dictionary.Dictionary)
		segmenter := dictionary.Sego()
		stopwords := dictionary.Stopwords()

		document.Init(segmenter, stopwords)
		document.LoadHTML(html)
	}
	
	return
}