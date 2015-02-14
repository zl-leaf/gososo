package analyzer

import(
	"github.com/huichen/sego"
	"math"
)

type Keyword struct {
	Text string
	Count int
	Token *sego.Token
	Document *Document
	Positions []*Position
}

type Position struct {
	Row int
	Index int
}

type Keywords []*Keyword

func (keywords Keywords) Len() int {return len(keywords)}
func (keywords Keywords) Less(i,j int) bool {
	if keywords[i].Weight() > keywords[j].Weight() {
		return true
	} else if keywords[i].Weight() == keywords[j].Weight() && keywords[i].Count > keywords[j].Count {
		return true
	}
	return false
} 
func (keywords Keywords) Swap(i,j int) {keywords[i],keywords[j] = keywords[j],keywords[i]}

func (keyword *Keyword) TF() float64{
	var f1 float64 = float64(keyword.Count)
	var f2 float64 = float64(keyword.Document.WordCount)
	return f1/f2
}

func (keyword *Keyword) IDF() float64{
	var f1 float64 = float64(keyword.Document.TotalFrequency())
	var f2 float64 = float64(keyword.Token.Frequency() + 1)
	return math.Log(f1/f2)
}

func (keyword *Keyword) PW() float64 {
	var res float64 = 0
	for _,pos := range keyword.Positions {
		res += 1/(5 * (float64(pos.Row) - 0.8) + float64(pos.Index))
	}
	return res
}

func (keyword *Keyword) Weight() float64 {
	return 1.5 * keyword.PW() + keyword.TF() * keyword.IDF()
}