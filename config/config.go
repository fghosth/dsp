package config

import (
	"runtime"
	"time"

	"github.com/RoaringBitmap/roaring"
)

var (
	RabbitMQURL = "amqp://derek:zaqwedcxs@mqm.nb.com:5672/"
	RedisURL    = "redis.nb.com:6380"
	// RedisURL                    = "localhost:6379"
	RedisPWD                    = ""
	RedisDB                     = 0
	CompaignURL                 = "http://54.178.104.115:5000/dsp-campaign/getMany"
	CPIDXNAME                   = "compaign.idx"                                    //索引本地缓存文件名
	RedisIndexKey               = "DSP:compaignIndex"                               //redis indexkey
	RedisCPVersionKey           = "DSP:CPIndexversion"                              //redis compaignversion key
	PerNum                      = 20                                                //初始化索引时每次读取几条数据
	Interval                    = 2                                                 //检查索引时间间隔 s
	Cashfix              uint32 = 1000000                                           //货币都要乘1000000
	MapLength                   = 10000                                             //设备，用户，位置等map预留大小
	Loglevel                    = "all"                                             //日志等级 all debug info  error
	NURL                        = "http://adx-smaato.newbidder.com/v1/ADXNotify"    //胜出通知链接。
	QPS                         = 1000000                                           //访问频率
	Timeout                     = 0                                                 //线程超时时间 单位s 0为永不超时
	Numcpu                      = runtime.NumCPU()                                  //多线程使用cpu数量
	NUMP                        = 200                                               //线程数
	CacheSize                   = 100 * 1024 * 1024                                 //缓存大小100MB
	CacheExpire                 = 600                                               //广告缓存保留时间600秒
	InfluxdbInsertNowURL        = "http://influx.newbidder.com/influx/v2/insertNow" //插入influxdb的链接,立刻插入
	InfluxdbInsertURL           = "http://influx.newbidder.com/influx/v2/insert"    //插入influxdb的链接，插入有缓存延迟
	InfluxdbSelectURL           = "http://influx.newbidder.com/influx/v2/select"    //插入influxdb的链接，插入有缓存延迟
	InfluxdbDatabase            = "dsp"                                             //投放广告缓存库
	InfluxdbTable               = "bidcache"                                        //投放广告缓存表 保存6小时数据
	TrackingEXName              = "dsp.tracking.notify"                             //成功后通知tracking的exname
	TrackingRouteKey            = "tracking"                                        //成功后通知tracking的routekey
	DSPEXName                   = "dsp.success.notify"                              //成功后通知dsp的exname
	DSPRouteKey                 = "success"                                         //成功后通知dsp的routekey
	Server                      = "dspSmaato1"                                      //服务器名称，每台必须不一样
	IndexEXName                 = "dsp.index.notify"                                //compaign更新后通知index的exname
	IndexRouteKey               = "index"                                           //compaign更新后通知index的routekey
	DSPextype                   = "fanout"                                          //rabbitmq extype 类型
	Bind                        = ":9900"                                           //服务端口号
	Maxconn                     = 100000                                            //最大连接数
	FilterTimeout               = 30                                                //用户筛选最大时间 ms
	CDNURL                      = "http://s3.us-east-2.amazonaws.com/"              //cdn地址
	ADX                         = "smaato"                                          //adx 服务商
	ADXReqTrackPath             = "adx.txt"                                         //adx请求日志路径
	ADXReqTrackISON             = false                                             //adx请求日志路径是否开启
	ADXReqTrackWin              = false                                             //竞价成功日志是否打开
	ADXReqTrackWinPath          = "win.log"                                         //竞价成功日志文件
	DSPSyncIndexServers         = []string{
		"http://localhost:9900/v1/SyncIndex",
		"http://localhost:9900/v1/SyncIndex",
	}

	DspTrackCIDs      []uint32                     //要跟踪问题的compaignid
	TrackCompaignISON bool     = true              //是否打开跟踪问题的compaignid
	TrackingFile               = "dsptracking.log" //用于追踪存放要跟踪问题的compaignid的记录
	//impressionURL后添加的参数
	IMPURL_ImageURLName = "imgurl"
	IMPURL_SiteORAPPID  = "WebsiteID"
	IMPURL_Domain       = "Domain"
	IMPURL_Cat          = "Cat"
	IMPURL_Page         = "Page"
	IMPURL_IDZone       = "IDZone"
	IMPURL_UA           = "ua"
	IMPURL_IP           = "ip"
	IMPURL_Country      = "country"
	IMPURL_City         = "city"
	IMPURL_Region       = "region"
	IMPURL_OS           = "os"
	IMPURL_Language     = "language"
	IMPURL_Model        = "model"
	IMPURL_OSV          = "osv"
	IMPURL_DeviceType   = "devicetype"
	IMPURL_DeviceID     = "deviceid"
	IMPURL_ADX          = "adx"
	IMPURL_ClickID      = "clickID"
	//===================以下配置文件不设置
	IsInit             = false            //是否是初始化启动 是则先删除所有消息
	CPMBidFix   uint32 = 1                //千次展示cpm要除1
	StartupTime        = time.Now().UTC() //程序启动时间

)

//adult list
var AdultList = []string{
	"adx-adsterra-add.newbidder.com",
	"adx-exads-adult.newbidder.com",
	"localhost3:9900",
}
var AdultBitmap = roaring.NewBitmap()

//adx编号
var ADXCode = map[string]uint32{
	"smaato":        1 << 0,  //9900
	"exads":         1 << 1,  //9910
	"mgid":          1 << 2,  //9920
	"popcash":       1 << 3,  //9901
	"propellerads":  1 << 4,  //9902
	"adsterra":      1 << 5,  //9903
	"zeropark":      1 << 6,  //9904
	"popundertotal": 1 << 7,  //9905
	"hilltopads":    1 << 8,  //9906
	"clickadu":      1 << 9,  //9907
	"eroadv":        1 << 10, //9908
}

//回掉连接参数字符串定义
// OID    string //广告位id
// Price uint32 //竞价价格
// CID   uint32 //compaignid
// UID   uint32 //用户id
// T     int64  //竞价成功时间
// Postion string //广告位置标识
// Device  string //设备标识
var NURLParam = map[string]string{
	"OID":     "oid",
	"Price":   "price",
	"CID":     "cid",
	"UID":     "uid",
	"Postion": "postion",
	"Device":  "device",
	"T":       "t",
	"User":    "user",
}

const (
	AUCTION_ID       = "${AUCTION_ID}"       // - ID of the bid request; from "Bid Request Object -> id" attribute.
	AUCTION_BID_ID   = "${AUCTION_BID_ID}"   //- ID of the bid; from "Bid Response Object -> bidid" attribute.
	AUCTION_IMP_ID   = "${AUCTION_IMP_ID}"   //- ID of the impression just won; from "Bid Request Object -> Impression Object -> id" attribute.
	AUCTION_SEAT_ID  = "${AUCTION_SEAT_ID}"  //- ID of the bidder seat for whom the bid was made; from "Bid Response Object -> Seat Bid Object -> Bid Object -> id" attribute.
	AUCTION_AD_ID    = "${AUCTION_AD_ID}"    //- ID of the ad markup the bidder wishes to serve; from "Bid Response Object -> Seat Bid Object -> Bid Object -> adid" attribute.
	AUCTION_PRICE    = "${AUCTION_PRICE}"    //- Settlement price using the same currency and units as the bid; from "Bid Response Object -> Seat Bid Object -> Bid Object -> price" attribute.
	AUCTION_CURRENCY = "${AUCTION_CURRENCY}" //- The currency used in the bid (explicit or implied); for confirmation only.
)

//过滤器编号
var FilterCode = map[string]uint32{
	"MaxBidFilter":          1,
	"DailyBudgetFilter":     1 << 1,
	"SpendStrategyFilter":   1 << 2,
	"TotalBudgetFilter":     1 << 3,
	"FreqCapFilter":         1 << 4,
	"BidderTimeFilter":      1 << 5,
	"DayPartingFilter":      1 << 6,
	"CountriesFilter":       1 << 7,
	"StatesFilter":          1 << 8,
	"CitiesFilter":          1 << 9,
	"ADXFilter":             1 << 10,
	"ControlListFilter":     1 << 11,
	"CategoriesFilter":      1 << 12,
	"CarriersFilter":        1 << 13,
	"OSFilter":              1 << 14,
	"AudiencesFilter":       1 << 15,
	"DailyPPBFilter":        1 << 16,
	"IPFilter":              1 << 17,
	"ConnectTypeFilter":     1 << 18,
	"DeviceTypeFilter":      1 << 19,
	"ImageTypeFilter":       1 << 20,
	"StatusFilter":          1 << 21,
	"NativeImageTypeFilter": 1 << 22,
	"AdultFilter":           1 << 23,
}
var ImageEMU = map[string]uint64{
	"300x250":  1,
	"320x50":   1 << 1,
	"320x480":  1 << 2,
	"728x90":   1 << 3,
	"768x1024": 1 << 4,
	"480x320":  1 << 5,
	"1024x768": 1 << 6,
}
var DSPSQL = map[string]string{
	"GetComByID":  "select * from BaseDspCampaign where id=",                                          //根据id获取compaign
	"GetTotalCom": "select count(id) as total from BaseDspCampaign where id in(2,81) and status=true", //获取总的compaign数量
	"GetCompaign": "select * from BaseDspCampaign where id in(2,81) and  status=true",                 //获取总的compaign数量id in(2,81,84,85,86,87,88,89,90,91)
}

//compaign状态
var CompaignStatus = map[string]uint8{
	CS_RUNNING: 1,
	CS_PENDING: 0,
}

const (
	CS_RUNNING = "running" //compaign状态 运行
	CS_PENDING = "pending" //compaign状态 暂停
	MQNADD     = "add"     //compaign改变时的action add
	MQNMODIFY  = "modify"  //compaign改变时的action modify
	MQNDEL     = "del"     //compaign改变时的action del
	MQNPENDING = "pending" //compaign改变时的action del
)
const (
	TEMPLATE_LEFT_TAG  = "{{"
	TEMPLATE_RIGHT_TAG = "}}"

	// SMAATO_TEMPLATE = "{\"seatbid\":[{\"bid\":[{\"nurl\":\"{{NURL}}\",\"crid\":\"{{CRID}}\",\"adomain\":[\"{{ADOMAIN}}\"],\"price\":{{PRICE}},\"id\":\"{{ID}}\",\"adm\":\"\",\"impid\":\"{{IMPID}}\",\"cid\":{{CID}}}]}],\"id\":\"{{ID}}\"}"
	BASE_BID_TEMPLATE = `{"id": "{{ID}}","cur": "USD","seatbid": [{"group": 1,"bid": [{"id": "{{ID}}","impid": "{{IMPID}}","price": {{PRICE}},"cid": "{{CID}}","adid": "{{ADID}}","adm": "{{ADM}}","nurl": "{{NURL}}"}]}]}`
	//=================================================SMAATO
	SMAATO_TEMPLATE = `{"seatbid":[{"bid":[{"nurl":"{{NURL}}","crid":"{{CRID}}","adomain":["{{ADOMAIN}}"],"price":{{PRICE}},"id":"{{ID}}","adm":"<ad xmlns:xsi=\"http:\/\/www.w3.org\/2001\/XMLSchema-instance\" xsi:noNamespaceSchemaLocation=\"smaato_ad_v0.9.xsd\" modelVersion=\"0.9\"><imageAd><clickUrl><![CDATA[{{CLICKURL}}]]><\/clickUrl><imgUrl><![CDATA[{{IMGURL}}]]><\/imgUrl><height>{{HEIGHT}}<\/height><width>{{WIDTH}}<\/width><beacons>{{BEACONURL}}<\/beacons><\/imageAd><\/ad>","impid":"{{IMPID}}","cid":{{CID}}}]}],"id":"{{ID}}"}`

	//=================================================EXADS
	EXADS_TEMPLATE            = `{"id": "{{ID}}","seatbid": [{"bid": [{"id": "{{ID}}","impid": "{{IMPID}}","price": {{PRICE}},"adm": "{{ADM}}","nurl": "{{NURL}}"}]}]}`
	EXADS_BANNER_XML_TEMPLATE = `<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><ad><imageAd><clickUrl><![CDATA[{{CLICKURL}}]]></clickUrl><imgUrl><![CDATA[{{IMGURL}}]]></imgUrl></imageAd></ad>`
	// EXADS_POPUP_XML_TEMPLATE   = `<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><ad><iframeAd><url><![CDATA[{{CLICKURL}}]]></url></iframeAd></ad>`
	EXADS_POPUP_XML_TEMPLATE   = `<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><ad><popunderAd><url><![CDATA[{{CLICKURL}}]]></url></popunderAd></ad>`
	EXADS_NATIVE_JSON_TEMPLATE = `{\"native\": {\"link\": {\"url\": \"{{CLICKURL}}\"},\"assets\": [{\"id\": 1,\"required\": 1,\"img\": {\"url\": \"{{IMGURL}}\",\"w\": {{WIDTH}},\"h\": {{HEIGHT}}}},{\"id\": 2,\"title\": {\"text\": \"{{TITLE}}\"}},{\"id\": 3,\"data\": {\"value\": \"{{DESCRIPTION}}\"}}]}}`

	//=================================================MGID
	// MGID_NATIVE_JSON_TEMPLATE = `{"native": {"ver": 1,"link": {"url": "{{CLICKURL}}"},"assets": [{"id": 1,"title": {"text": "{{TITLE}}"}},{"id": 2,"data": {"value": ""},"img": {"url": "{{IMGURL}}","w": {{WIDTH}},"h": {{HEIGHT}}}},{"id": 3,"data": {"value": "{{DESCRIPTION}}"}}]}}`
	MGID_NATIVE_JSON_TEMPLATE = `{\"native\": {\"ver\": 1,\"link\": {\"url\": \"{{CLICKURL}}\"},\"assets\": [{\"id\": 1,\"title\": {\"text\": \"{{TITLE}}\"}},{\"id\": 2,\"data\": {\"value\": \"\"},\"img\": {\"url\": \"{{IMGURL}}\",\"w\": {{WIDTH}},\"h\": {{HEIGHT}}}},{\"id\": 3,\"data\": {\"value\": \"{{DESCRIPTION}}\"}}]}}`

	//=================================================POPCASH
	POPCASH_POPUP_JSON_TEMPLATE = `<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><ad><popunderAd><url><![CDATA[{{CLICKURL}}]]></url></popunderAd></ad>`
	//=================================================PROPELLERADS
	PROPELLERADS_POPUP_JSON_TEMPLATE = `{"id": "{{ID}}","seatbid" :[{"bid": [{"id": "{{ID}}","impid": "popup","ext": {"url" : "{{CLICKURL}}"},"price": {{PRICE}}}]}]}`
	//=================================================Adsterra
	ADSTERRA_POPUP_JSON_TEMPLATE = `{"id": "{{ID}}","bidid": "{{ID}}","seatbid": [{"bid": [{"id": "1","impid": "{{IMPID}}","price": 0.1,"adid": "{{ID}}","ext": {"url": "{{CLICKURL}}"}}]}],"cur": "USD"}`
	//=================================================Zeropark
	ZEROPARK_POPUP_XML_TEMPLATE = `<?xml version=\"1.0\" encoding=\"UTF8\"?><results><listing><url><![CDATA[{{CLICKURL}}]]</url><bid><![CDATA[{{PRICE}}]]</bid><listing></results>`
	//=================================================Popundertotal
	POPUNDERTOTAL_POPUP_XML_TEMPLATE = ``
	//=================================================hilltopads
	HILLTOPADS_POPUP_JSON_TEMPLATE = `{"id": "{{ID}}","bidid": "{{ID}}","cur": "USD","seatbid": [{"bid": [{"id": "{{ID}}","impid": "1","price": {{PRICE}},"nurl": "{{NURL}}","adid": "{{ID}}","adomain": ["{{ADOMAIN}}"],"adm": "{{CLICKURL}}","cid": "{{CID}}","crid": "{{CRID}}"}]}]}`
	//=================================================eroadv
	Eroadv_BANNER_JSON_TEMPLATE = `<iframe scrolling=\"no\" width=\"{{WIDTH}}\" height=\"{{HEIGHT}}\" frameborder=\"0\" src=\"{{IMGURL}}\" marginwidth=\"0\" marginheight=\"0\" allowtransparency=\"true\"><\/iframe>`
	Eroadv_POP_JSON_TEMPLATE    = `{{CLICKURL}}`
)
