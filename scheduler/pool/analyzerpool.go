package pool
import(
    "regexp"
    "strings"
    "errors"
)

type AnalyzerPool struct {
    specialPoolMap map[string]*Pool
    allPool *Pool
    hasAll bool
}

func NewAnalyzerPool() *AnalyzerPool {
    analyzerPool := &AnalyzerPool{}
    analyzerPool.specialPoolMap = make(map[string]*Pool)
    analyzerPool.allPool = NewPool()
    analyzerPool.hasAll = false
    return analyzerPool
}

func (analyzerPool *AnalyzerPool) Add(regexpString string, e interface{}) {
    if strings.ToLower(regexpString) == "all" {
        analyzerPool.hasAll = true
        analyzerPool.allPool.Add(e)
    } else {
        specialPool,exist := analyzerPool.specialPoolMap[regexpString]
        if !exist {
            specialPool = NewPool()
            analyzerPool.specialPoolMap[regexpString] = specialPool
        }
        specialPool.Add(e)
    }
}

func (analyzerPool *AnalyzerPool) Get(u string) (value interface{}, err error) {
    for regexpString,p := range analyzerPool.specialPoolMap {
        r := regexp.MustCompile(regexpString)
        if r.Match([]byte(u)) {
            value = p.Get()
            return
        }
    }
    if analyzerPool.hasAll {
        value = analyzerPool.allPool.Get()
        return
    } else {
        err = errors.New("没有匹配的分析器")
        return
    }
}
