package pool
import(
	"time"
	"../../utils/queue"
)

type Pool struct {
	elements *queue.Queue
}

func NewDownloaderPool() (pool *Pool){
	pool = &Pool{}
	pool.elements = queue.New()
	return
}

func (pool *Pool)Add(e interface{}) {
	pool.elements.Add(e)
}

func (pool *Pool)Get() (value interface{}) {
	for {
		if !pool.elements.Empty() {
			e,_ := pool.elements.Head()
			value =  e.Value
			return
		}
		time.Sleep(1 * time.Second)
	}
}