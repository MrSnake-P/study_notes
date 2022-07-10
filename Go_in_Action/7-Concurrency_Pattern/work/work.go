// 使用无缓冲的通道来创建一个goroutine池
// 这些goroutine执行并控制一组工作，让其并发执行
//
// 无缓冲的通道保持两个goroutine之间的数据交换
// 这种使用方法允许使用者直到什么时候goroutine池
// 正在执行工作，而且如果池里的所有goroutine都忙
// 无法接受新的工作的时候，也能及时通过通道来通知调用者

package work

import "sync"

type Worker interface {
	Task()
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			// 直到收到一个Worker接口直，开始执行方法
			// 通道被关闭时候，循环结束
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *Pool) Run(w Worker) {
	p.work <- w
}

func (p *Pool) Shutdown() {
	// 关闭通过，等待所有goroutine完成
	close(p.work)
	p.wg.Wait()
}
