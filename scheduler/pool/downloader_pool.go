package pool
import(
	"net"
	"time"
	"../../utils/queue"
)

type DownloaderPool struct {
	elements *queue.Queue
}

func NewDownloaderPool() (pool *DownloaderPool){
	pool = &DownloaderPool{}
	pool.elements = queue.New()
	return
}

func (pool *DownloaderPool)Add(conn net.Conn) {
	pool.elements.Add(conn)
}

func (pool *DownloaderPool)Get() (conn net.Conn) {
	result := make(chan net.Conn)
	go func() {
		for {
			if !pool.elements.Empty() {
				e,_ := pool.elements.Head()
				result <- e.Value.(net.Conn)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return <- result
}