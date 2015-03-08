package analyzer

import(
	"regexp"
	"strings"
	"github.com/huichen/sego"
	"sort"
)

type Document struct {
	Title string
	MainContent string
	WordCount int
	Keywords Keywords
	Segmenter sego.Segmenter
	Stopwords map[string]int
}

func (doc *Document) Init(segmenter sego.Segmenter, stopwords map[string]int) {
	doc.Segmenter = segmenter
	doc.Stopwords = stopwords
}

func (doc *Document)LoadHTML(html string) {
	doc.Title, doc.MainContent = getMainContent(html)
	words := doc.Words()
	doc.WordCount = len(words)
	doc.Keywords = getKeywords(doc, words)
}

func (doc *Document) Words() ([]*sego.Token){
	segmenter := doc.Segmenter

    // 分词
    text := []byte(doc.Title + "\n" + doc.MainContent)
    segments := segmenter.Segment(text)

    words := make([]*sego.Token, 0)
    for _, segment := range segments {
    	words = append(words, segment.Token())
    }
    return words
}

func (doc *Document) TotalFrequency() int64 {
	return doc.Segmenter.Dictionary().TotalFrequency()
}

func getMainContent(html string) (title string, mainCount string) {
	title = getHTMLTitle(html)

	hrefRegexp := regexp.MustCompile(`(<(head|script|style|noscript).*?>[\s\S]*?<\/(head|script|style|noscript)>|<[^>]+>|&nbsp)`)
	html = hrefRegexp.ReplaceAllString(html, "")
	lines := strings.Split(html, "\n")
	charCount := len(html)
	lineCount := len(lines)

	startPos := -1
	endPos := -1

	flag2 := false// 是否搜索完成正文
	for i:=0; i<lineCount; i+=5 {
		tmpCharCount := 0
		tmpLineCount := 0
		tmpStr := ""
		for j:=0;j<5 && i+j<lineCount;j++ {
			line := strings.TrimSpace(lines[i+j])

			if strings.TrimSpace(line) != "" {
				tmpCharCount += len(strings.Replace(line, " ", "", -1))
				tmpLineCount ++
				tmpStr += line
			}
		}


		f1 := float32(tmpCharCount)
		f2 := float32(charCount)
		f := f1/f2

		if f > 0.02 {
			flag2 = true
			if startPos < 0 {
				// 搜索开头
				flag := false// 判断是否搜索到非空行
				emptyLine := 0
				for j:=4;j>=0;j-- {
					if i+j >= lineCount {
						continue
					}

					line := strings.TrimSpace(lines[i+j])
					if line != "" {
						startPos = i + j
						flag = true
					}
					if flag {
						if line == "" {
							emptyLine ++
						} else {
							emptyLine = 0
						}
					}
					if emptyLine >= 2 {
						break
					}
				}
			}
		} else {
			// 搜索结尾
			if flag2 {
				flag2 = false
				endPos = i
				emptyLine := 0
				for j:=0;j<5 && i+j<lineCount;j++ {
					line := strings.TrimSpace(lines[i+j])
					if strings.TrimSpace(line) == "" {
						emptyLine++
					} else {
						endPos = i + j
						emptyLine = 0
					}

					if emptyLine >= 2 {
						break
					}
				}
			}
		}
	}

	mainCount = ""
	for i,line := range lines {
		if i >= startPos && (i <= endPos || endPos == -1) {
			line := strings.TrimSpace(line)
			if line != "" {
				mainCount += line + "\n"
			}
		}
	}

	return title, mainCount
}

func getHTMLTitle(html string) string{
	title := ""
	var hrefRegexp = regexp.MustCompile(`<title.*?>(.*?[^<])</title>`)
	match := hrefRegexp.FindAllStringSubmatch(html, -1)
	if match != nil {
		for _,m := range match {
	 		title += m[1] + " "
	 	}
	}
	return title
}

func getKeywords(doc *Document, tokens []*sego.Token) Keywords{
	keyWords := Keywords{}

	stopwords := doc.Stopwords
	wordMap := make(map[string]*Keyword)

	row := 1
	index := 1
	for _,token := range tokens {
		text := token.Text()
		if text == "\n" {
			if index != 1 {
				row ++
			}
			index = 1
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		if _,isStoped := stopwords[text];!isStoped {
			keyWord,exist := wordMap[text]
			if !exist {
				keyWord = &Keyword{Text:text, Count:0, Token:token, Document:doc}
				wordMap[text] = keyWord
			}
			keyWord.Positions = append(keyWord.Positions, &Position{Row:row, Index:index})
			keyWord.Count ++
			index ++
		}
	}

	for _,keyWord := range wordMap {
		keyWords = append(keyWords, keyWord)
	}

	sort.Sort(keyWords)
	if len(keyWords) > 30 {
		return keyWords[:30]
	} else {
		return keyWords
	}
}
