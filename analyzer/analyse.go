package analyzer

import(
	"os"

	"github.com/zl-leaf/gososo/utils/dictionary"
	"github.com/zl-leaf/extract"
	"github.com/zl-leaf/extract/exp"
)

func (analyzer *Analyzer) analyse(f string) (document *exp.Document, err error) {
	file, err := os.Open(f)
	if err != nil {
		return
	}
	defer file.Close()

	// document = &Document{}
	component,exist := analyzer.context.GetComponent("dictionary")
	if exist {
		dictionary := component.(*dictionary.Dictionary)
		segmenter := dictionary.Sego()
		stopwords := dictionary.Stopwords()

		// document.Init(segmenter, stopwords)
		// document.LoadHTML(html)
		extractor := extract.New(segmenter,stopwords)
		document = extractor.Extract(file)
	}

	return
}
