package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"

	"github.com/fsnotify/fsnotify"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/k0kubun/pp"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"golang.org/x/net/netutil"
	"jvole.com/dsp/config"
	"jvole.com/dsp/grpc"
	"jvole.com/dsp/index"
	"jvole.com/dsp/kafka"
	"jvole.com/dsp/rabbitmq"
	"jvole.com/dsp/service"
	"jvole.com/dsp/util"
)

type configmap struct {
	RabbitMQURL          string
	RedisURL             string
	RedisPWD             string
	RedisDB              int
	CompaignURL          string
	CPIDXNAME            string //索引本地缓存文件名
	RedisIndexKey        string //redis indexkey
	RedisCPVersionKey    string //redis compaignversion key
	PerNum               int    //初始化索引时每次读取几条数据
	Interval             int    //检查索引时间间隔 s
	Cashfix              uint32 //货币都要乘1000000
	MapLength            int    //设备，用户，位置等map预留大小
	DSPSQL               map[string]string
	NURL                 string //胜出通知链接。
	QPS                  int
	ADXCode              map[string]uint32
	Timeout              int      //线程超时时间 单位s 0为永不超时
	Numcpu               int      //多线程使用cpu数量
	NUMP                 int      //线程数
	CacheSize            int      //缓存大小100MB
	CacheExpire          int      //缓存过期时间 s
	InfluxdbInsertNowURL string   //插入influxdb的链接,立刻插入
	InfluxdbInsertURL    string   //插入influxdb的链接，插入有缓存延迟
	InfluxdbSelectURL    string   //插入influxdb的链接，插入有缓存延迟
	InfluxdbDatabase     string   //投放广告缓存库
	InfluxdbTable        string   //投放广告缓存表 保存6小时数据
	TrackingEXName       string   //成功后通知tracking的exname
	TrackingRouteKey     string   //成功后通知tracking的routekey
	DSPEXName            string   //成功后通知tracking的exname
	DSPRouteKey          string   //成功后通知tracking的routekey
	Server               string   //服务器名称，每台必须不一样
	IndexEXName          string   //compaign更新后通知index的exname
	IndexRouteKey        string   //compaign更新后通知index的routekey
	DSPextype            string   //rabbitmq extype 类型
	Bind                 string   //服务端口号
	Maxconn              int      //最大连接数
	FilterTimeout        int      //用户筛选最大时间 ms
	CDNURL               string   //cdnurl
	ADX                  string   //adx
	ADXReqTrackPath      string   //adx请求日志路径
	ADXReqTrackISON      bool     //adx请求日志路径是否开启
	Loglevel             string   //日志等级
	ADXReqTrackWin       bool     //竞价成功日志是否打开
	ADXReqTrackWinPath   string   //竞价成功日志文件
	DSPSyncIndexServers  []string //dsp同步索引服务列表
	DspTrackCIDs         []uint32 //要跟踪问题的compaignid
	TrackCompaignISON    bool     //是否打开跟踪问题的compaignid
	TrackingFile         string   //用于追踪存放要跟踪问题的compaignid的记录
	KafkaURL             string   //kafka连接
	GRPCPEM              string
	GRPCKEY              string
	GrpcAddress          string
}

var (
	cfg         *configmap
	defaultFile = "server.yaml"
	Viper       = viper.New()
)

func main() {
	file := flag.String("f", "./"+defaultFile, "config file path")
	flag.Parse()
	exist, _ := util.PathExists(*file)
	if exist {
		Viper.SetConfigType("yaml")
		// util.Viper.SetConfigName(".cfg")
		Viper.AddConfigPath(".")
		loadConfig(*file)
		//监控配置文件变化
		Viper.WatchConfig()
		Viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("配置文件变更，重新生效")
			loadConfig(*file)
			pp.Println(cfg)
		})
	} else {
		log.Printf("找不到配置文件，请加参数,eg：-f /etc/server.yaml")
		os.Exit(0)
	}

	// pp.Println(cfg)
	// for k, v := range config.FilterCode {
	// 	fmt.Println(k, v)
	// }
	//初始化标识
	config.IsInit = false
	//=================================初始化index
	log.Println("初始化索引...")
	index.CPINDEX.SetupIndex()
	// index.CPINDEX.SaveRedis()
	index.CPINDEX.IndexCheck() //定时任务检查，根据index
	// index.CPINDEX.SaveDisk()
	// pp.Println("comp", len(index.CPINDEX.Compaign), index.CPINDEX.Bitmap.TypeBanner.GetCardinality())
	//==================================初始化kafka producer
	log.Println("初始化Kafka Producer...")
	kafka.KafKaProducer = kafka.NewDspMsg()
	//==================================初始化rabbitMQ
	log.Println("初始化rabbitMQ...")
	conf := &rabbitmq.Mqserver{
		Url:    config.RabbitMQURL,
		Online: true,
	}
	//=================================grpc服务
	log.Println("初始化GRPC...")
	go grpc.Run()
	//index message
	MQIndexConn := rabbitmq.NewRabbitMQ(config.IndexEXName, config.DSPextype, *conf)
	indexreceiver := rabbitmq.NewIndexReceiver("index_"+config.Server, config.IndexRouteKey)
	MQIndexConn.RegisterReceiver(indexreceiver)
	go MQIndexConn.Start()
	//success message
	MQSuccessConn := rabbitmq.NewRabbitMQ(config.DSPEXName, config.DSPextype, *conf)
	rabbitmq.RabbitMQConn = MQSuccessConn
	for i := 0; i < 5; i++ {
		successreceiver := rabbitmq.NewSuccessReceiver("success_"+config.Server, config.DSPRouteKey)
		MQSuccessConn.RegisterReceiver(successreceiver)
	}
	go MQSuccessConn.Start()

	//========http服务
	fieldKeys := []string{"method"}

	// logging := kitlog.With(util.KitLogger, "component", "DSP")
	var bs service.DSPBidder
	bs = service.GetService(config.ADXCode[config.ADX], &util.KitLogger) //获取adx服务
	pp.Println("========ADX:======", reflect.TypeOf(bs).String())
	bs = service.NewLoggingService(&util.KitLogger, bs)
	bs = service.NewInstrumentingService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "DSP_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "DSP_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		bs)
	httpLogger := kitlog.With(util.KitLogger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/v1/", service.MakeHandler(bs, httpLogger))
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	l, err := net.Listen("tcp", config.Bind)
	if err != nil {
		pp.Println(err)
	}

	defer l.Close()
	// pp.Println(config.Maxconn)
	l = netutil.LimitListener(l, config.Maxconn)
	log.Println("http服务启动成功...")
	http.Serve(l, nil)

	//TODO 平滑重启 安全退出  官方的访问限制方案对用户并不友好，可以i自己写
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

//读取配置文件
func loadConfig(file string) {

	// util.Viper.AddConfigPath("/Users/derek/project/go/src/jvole.com/monitor/")
	Viper.SetConfigFile(file)

	err := Viper.ReadInConfig() // 读取配置文件
	if err != nil {             // 加载配置文件错误
		log.Printf("配置文件加载错误")
	}
	cfg = &configmap{}
	err = Viper.Unmarshal(cfg)
	if err == nil { //覆盖默认配置
		config.RabbitMQURL = cfg.RabbitMQURL
		config.RedisURL = cfg.RedisURL
		config.RedisPWD = cfg.RedisPWD
		config.RedisDB = cfg.RedisDB
		config.CompaignURL = cfg.CompaignURL
		config.CPIDXNAME = cfg.CPIDXNAME
		config.RedisIndexKey = cfg.RedisIndexKey
		config.RedisCPVersionKey = cfg.RedisCPVersionKey
		config.PerNum = cfg.PerNum
		config.Interval = cfg.Interval
		config.Cashfix = cfg.Cashfix
		config.MapLength = cfg.MapLength
		config.DSPSQL = cfg.DSPSQL
		config.NURL = cfg.NURL
		// config.NURL = util.DealXMLURL(cfg.NURL)
		config.QPS = cfg.QPS
		config.ADXCode = cfg.ADXCode
		config.Numcpu = cfg.Numcpu
		config.NUMP = cfg.NUMP
		config.Timeout = cfg.Timeout
		config.CacheSize = cfg.CacheSize
		config.CacheExpire = cfg.CacheExpire
		config.InfluxdbInsertNowURL = cfg.InfluxdbInsertNowURL
		config.InfluxdbInsertURL = cfg.InfluxdbInsertURL
		config.InfluxdbSelectURL = cfg.InfluxdbSelectURL
		config.InfluxdbDatabase = cfg.InfluxdbDatabase
		config.InfluxdbTable = cfg.InfluxdbTable
		config.TrackingEXName = cfg.TrackingEXName
		config.TrackingRouteKey = cfg.TrackingRouteKey
		config.DSPEXName = cfg.DSPEXName
		config.DSPRouteKey = cfg.DSPRouteKey
		config.Server = cfg.Server
		config.IndexEXName = cfg.IndexEXName
		config.IndexRouteKey = cfg.IndexRouteKey
		config.DSPextype = cfg.DSPextype
		config.Bind = cfg.Bind
		config.Maxconn = cfg.Maxconn
		config.FilterTimeout = cfg.FilterTimeout
		config.CDNURL = cfg.CDNURL
		config.ADX = cfg.ADX
		config.ADXReqTrackPath = cfg.ADXReqTrackPath
		config.ADXReqTrackISON = cfg.ADXReqTrackISON
		config.Loglevel = cfg.Loglevel
		config.ADXReqTrackWin = cfg.ADXReqTrackWin
		config.ADXReqTrackWinPath = cfg.ADXReqTrackWinPath
		config.DSPSyncIndexServers = cfg.DSPSyncIndexServers
		config.DspTrackCIDs = cfg.DspTrackCIDs
		config.TrackingFile = cfg.TrackingFile
		config.KafkaURL = cfg.KafkaURL
		config.GRPCKEY = cfg.GRPCKEY
		config.GRPCPEM = cfg.GRPCPEM
		config.GrpcAddress = cfg.GrpcAddress
	} else {
		log.Println("配置文件错误:", err)
		os.Exit(0)
	}

	//===========log 设置
	// serv := &loginflux.InfluxServer{}
	// lf := loginflux.NewLoginflux(*serv)

	// KitLogger = log.NewJSONLogger(log.NewSyncWriter(lf))
	util.KitLogger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	switch config.Loglevel {
	case "all":
		util.KitLogger = level.NewFilter(util.KitLogger, level.AllowAll())
	case "debug":
		util.KitLogger = level.NewFilter(util.KitLogger, level.AllowDebug())
	case "info":
		util.KitLogger = level.NewFilter(util.KitLogger, level.AllowInfo())
	case "error":
		util.KitLogger = level.NewFilter(util.KitLogger, level.AllowError())
	}
	util.KitLogger = kitlog.With(util.KitLogger, "ts", kitlog.DefaultTimestampUTC)

}
