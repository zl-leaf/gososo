package dictionary
import(
	"os"
	"bufio"
	"strings"
	"io"

	"github.com/huichen/sego"
)

type Dictionary struct {
	segmenter sego.Segmenter
	stopwords map[string]int
}

func New(dictionaryPath, stopwordsPath string) (dict *Dictionary) {
	dict = &Dictionary{}
	dict.segmenter.LoadDictionary(dictionaryPath)
	stopwords,err := getStopwrods(stopwordsPath)
	if err == nil {
		dict.stopwords = stopwords
	} else {
		dict.stopwords = make(map[string]int)
	}
	return dict
}

func (dict *Dictionary) Sego() sego.Segmenter {
	return dict.segmenter
}

func (dict *Dictionary) Stopwords() map[string]int {
	return dict.stopwords
}

func getStopwrods(f string) (map[string]int,error){
	stopwords := make(map[string]int)
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		word, err := rd.ReadString('\n')
		word = strings.TrimSpace(word)
        if io.EOF == err {
            break
        }
		if err != nil {
            return nil, err
        }
        stopwords[word] = 0
	}
	return stopwords, nil
}