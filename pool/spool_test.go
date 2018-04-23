package pool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"jvole.com/dsp/config"
)

func TestSpool_NewSpool(t *testing.T) {
	sp := NewSpool(10, 100)
	sp.SetFinishCallback(finish)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		var result int
		worksp := &Worksp{
			result,
			func() interface{} {
				return job(r.Intn(10000000))
			},
			dealres,
		}

		sp.AddTask(*worksp)
	}
	go func() { //超时中断任务
		time.Sleep(time.Second * time.Duration(config.FilterTimeout))
		// sp.Stop()
	}()
	sp.Start()

}

func finish() {
	fmt.Println("任务结束")
}

func job(n int) interface{} {
	time.Sleep(time.Second * 1)
	return n * 2
}

func dealres(obj interface{}) (err error) {
	if n, ok := obj.(int); ok {
		fmt.Println(n)
	}
	return
}
