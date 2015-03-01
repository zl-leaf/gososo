package msg

type SearchResultMsg struct {
	Result int `json:"result"`
	Msg string `json:"msg"`
	Data SearchDatas `json:"data"`
}

type SearchResultObj struct {
	URL string `json:"url"`
	Title string `json:"title"`
	Description string `json:"description"`
	Keywords string `json:"keywords"`
	Weight float64 `json:"weight"`
}

type SearchDatas []*SearchResultObj

func (searchDatas SearchDatas) Len() int {return len(searchDatas)}
func (searchDatas SearchDatas) Less(i,j int) bool {
	if searchDatas[i].Weight > searchDatas[j].Weight {
		return true
	} 
	return false
} 
func (searchDatas SearchDatas) Swap(i,j int) {searchDatas[i],searchDatas[j] = searchDatas[j],searchDatas[i]}