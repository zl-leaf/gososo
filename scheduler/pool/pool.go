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

func (pool *Pool)Get() (e interface{}) {
	result := make(chan interface{})
	go func() {
		for {
			if !pool.elements.Empty() {
				e,_ := pool.elements.Head()
				result <- e.Value
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return <- result
}