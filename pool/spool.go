package pool

import (
	"log"
	"path"
	"runtime"

	"github.com/go-kit/kit/log/level"
)

type Worksp struct {
	Result  interface{}
	JobFunc func() interface{}
	ResFunc func(interface{}) error
}

type Spool struct {
	Queue         chan Worksp
	RuntineNumber int //线程数
	Total         int //任务数量

	Result         chan Worksp
	FinishCallback func()
}

//初始化
func NewSpool(runtineNumber int, total int) *Spool {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置所有cpu核心可用
	spool := &Spool{}
	spool.RuntineNumber = runtineNumber
	spool.Total = total
	spool.Queue = make(chan Worksp, total)
	spool.Result = make(chan Worksp, total)
	return spool
}

func (sp *Spool) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("pool意外崩溃", err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	//开启 number 个goruntine
	for i := 0; i < sp.RuntineNumber; i++ {

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("chan已关闭，任务超时", err) // 这里的err其实就是panic传入的内容，55
				}
			}()
			for {
				task, ok := <-sp.Queue
				if !ok {
					break
				}
				task.Result = task.JobFunc()

				sp.Result <- task
			}
		}()
	}

	//获取每个任务的处理结果
	for j := 0; j < sp.RuntineNumber; j++ {
		res, ok := <-sp.Result
		if !ok {
			break
		}
		err := res.ResFunc(res.Result)
		if err != nil {
			pc, file, line, _ := runtime.Caller(1)
			f := runtime.FuncForPC(pc)
			level.Error(*logger).Log(
				"method", f.Name(),
				"file", path.Base(file),
				"line", line,
				"msg", "返回函数结果错误",
				"err", err,
			)
			break
		}

	}

	//结束回调函数
	if sp.FinishCallback != nil {
		sp.FinishCallback()
	}
}

//关闭
func (sp *Spool) Stop() {
	close(sp.Queue)
	close(sp.Result)
}

func (sp *Spool) AddTask(task Worksp) {
	sp.Queue <- task
}

func (sp *Spool) SetFinishCallback(fun func()) {
	sp.FinishCallback = fun
}
