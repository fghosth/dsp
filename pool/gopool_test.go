package pool

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

func aTest_workerPool_AddWorkers(t *testing.T) {
	shutdown := make(chan int)
	timeout := 0
	numcpu := runtime.NumCPU()
	jobs := make(chan interface{})
	result := make(chan interface{})
	poolset := &Setting{
		shutdown,
		timeout,
		5,
		numcpu,
	}
	gopool := Newgopool(*poolset, jobs)
	gopool.AddWorkers(func(job interface{}) {
		works(result, job)
	})
	go func() {
		for i := 0; i < 50; i++ {
			jobs <- i
		}
	}()

	go func() {
		for {
			select {
			case res, ok := <-result:
				if ok {
					fmt.Println(res)
				} else {
					// fmt.Println("channel closed")
				}
			}
		}
	}()
	time.Sleep(time.Second * 2)
	log.Println(gopool.GetGoNum())
	// for i := 0; i < 5; i++ {
	// 	// shutdown <- 0
	// 	// timeout <- 0
	// 	destory <- 0
	// }
	// gopool.DestoryAll()
	time.Sleep(time.Second * 4)
	log.Println(gopool.GetGoNum())

	time.Sleep(time.Second * 10)
	close(shutdown)
	close(jobs)
	close(result)
}

func works(result chan interface{}, job interface{}) {
	// time.Sleep(time.Second * 2)
	if st, ok := job.(int); ok {
		result <- st * 2
	}

}
