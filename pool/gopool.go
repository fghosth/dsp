package pool

import (
	"path"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

func init() {
	logger = &util.KitLogger
	log.With(*logger, "component", "RabbitMQ")

}

var (
	logger   *log.Logger
	Jobs     = make(chan interface{}, config.NUMP)
	Result   = make(chan interface{}, config.NUMP)
	Shutdown = make(chan int)
)
var WorkpoolConn = Newgopool(*&Setting{Shutdown, config.Timeout, config.NUMP, config.Numcpu}, Jobs)

type WorkerPool interface {
	/*
			  添加任务
		    @param work 添加的任务
		    @param num  次任务的协程数量
		    @return error
	*/
	AddWorkers(work func(job interface{})) error
	/*获得正在运行的线程数
	@return int32 数量
	*/
	GetGoNum() int32
	/*销毁所有进程
	 */
	DestoryAll()
}

type workerPool struct {
	shutdown chan int         //关机信号
	timeout  chan int         //超时信号
	time     int              //超时的时间，单位秒
	destory  chan int         //销毁信号
	job      chan interface{} //任务
	num      int              //线程数
	numcpu   int              //使用多核工作的cpu数量
	gonum    int32            //正在运行协程的数量
	count    int              //计数器timeout用
}

type Setting struct {
	Shutdown chan int //关机信号
	Timeout  int      //超时时间,s,0为永不超时
	Num      int      //线程数
	Numcpu   int      //使用多核工作的cpu数量
}

var interval = 1

func Newgopool(set Setting, job chan interface{}) *workerPool {
	destorysign := make(chan int)
	timeoutsign := make(chan int)
	wp := &workerPool{
		shutdown: set.Shutdown,
		timeout:  timeoutsign,
		time:     set.Timeout,
		destory:  destorysign,
		job:      job,
		num:      set.Num,
		numcpu:   set.Numcpu,
		gonum:    0,
		count:    0,
	}
	//超时处理的协程
	go func() {
		// log.Printf("超时监控线程开始")
		// defer log.Printf("超时监控线程结束")
		if wp.time == 0 {
			return
		}
		for range time.Tick(time.Duration(interval) * time.Second) { //每秒检查一次
			wp.count++
			if wp.count > wp.time { //超时退出所有协程
				for i := 0; i < int(wp.GetGoNum()); i++ {
					wp.timeout <- 0
				}
			}
		}
	}()
	return wp
}

func (wp *workerPool) AddWorkers(work func(job interface{})) error {
	var err error
	runtime.GOMAXPROCS(wp.numcpu) //设置cpu的核的数量，从而实现高并发
	for i := 0; i < wp.num; i++ {
		// time.Sleep(time.Nanosecond * 10000)
		go func() {
			atomic.AddInt32(&wp.gonum, 1)
			// log.Printf("线程%d开始", i)
			// defer log.Printf("线程%d结束", i)
			defer func() {
				atomic.AddInt32(&wp.gonum, -1)
			}()
			for {
				select {
				case job, ok := <-wp.job:
					if ok {
						work(job)
					}
				case _, ok := <-wp.destory:
					if ok {
						pc, file, line, _ := runtime.Caller(1)
						f := runtime.FuncForPC(pc)
						level.Info(*logger).Log(
							"method", f.Name(),
							"file", path.Base(file),
							"line", line,
							"msg", "主动销毁导致线程退出",
						)
						return
					}
				case _, ok := <-wp.timeout:
					if ok {
						pc, file, line, _ := runtime.Caller(1)
						f := runtime.FuncForPC(pc)
						level.Info(*logger).Log(
							"method", f.Name(),
							"file", path.Base(file),
							"line", line,
							"msg", "超时导致线程退出",
						)
						return
					}
				case _, ok := <-wp.shutdown:
					if ok {
						pc, file, line, _ := runtime.Caller(1)
						f := runtime.FuncForPC(pc)
						level.Info(*logger).Log(
							"method", f.Name(),
							"file", path.Base(file),
							"line", line,
							"msg", "关闭程序导致线程退出",
						)
						return
					}
				}
			}
		}()

	}
	return err
}

func (wp *workerPool) GetGoNum() int32 {
	return atomic.LoadInt32(&wp.gonum)
}
func (wp *workerPool) DestoryAll() {
	for i := 0; i < int(atomic.LoadInt32(&wp.gonum)); i++ {
		wp.destory <- 0
	}
	close(wp.destory)
}
