package pool

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	maxGoroutines   = 25 // 要使用的goroutine的数量
	pooledResources = 2  // 池中的资源的数量
)

// 分配一个独一无二的id
var idCounter int32

func TestPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	p, err := New(createConnection, pooledResources)
	if err != nil {
		t.Fatal(err)
	}

	// 使用池里的连接来完成查询
	for query := 0; query <= maxGoroutines; query++ {
		// 每个goroutine需要自己复制一份查询值的副本，不然所有查询会共享同一个查询变量
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}
	wg.Wait()
	t.Log("Shutdown Program.")
	p.Close()
}

// dbConnection模拟要共享的资源
type dbConnection struct {
	Id int32
}

func (db *dbConnection) Close() error {
	log.Println("Close: Connection", db.Id)
	return nil
}

func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create: New Connection", id)

	return &dbConnection{id}, nil
}

func performQueries(query int, p *Pool) {
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Release(conn)

	time.Sleep(time.Duration(rand.Intn(199)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).Id)
}
