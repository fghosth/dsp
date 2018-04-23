package model

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/buger/jsonparser"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

var (
	logger *log.Logger
)

func init() {
	logger = &util.KitLogger
	log.With(*logger, "component", "Model")
}

//compaign 的结构
type Compaign struct {
	ID                      uint32
	UID                     uint32            //用户id
	TimeZone                TZ                //时区
	Type                    uint8             //广告类型NATIVE：1,BANNER：2,POPUP：3
	Maxbid                  uint32            //最大千次展示金额。千次内不超过max floorbid 比较
	Status                  uint8             //状态:runing 1,pending 0
	UnlimitBudget           bool              //是否不限制预算 true:不限制 false限制
	DailyBudget             uint32            //日预算
	DailyBudgetRecores      DailyBudgetRecore //日预算记录
	SpendStrategy           uint8             //1为ASAP，2为smooth
	DailyPerPlacementBudget uint32            //每天每个广告位的预算
	DailyPPBRecords         DailyPPBRecord    //每天每个广告位的预算记录
	TotalBudget             uint32            //总预算
	TotalBudgetRecords      TotalBudgetRecord //总预算记录
	FreqCapEnabled          bool              //是否开启投放频率
	FreqCapType             uint8             //DEVICE-设备:1，USER-用户:2
	FreqCountLimit          uint32            //投放次数限制
	FreqTimeWindow          uint8             //3,6,12,24小时
	FreqRecords             FreqRecord        //投放记录
	StartDate               time.Time         //投放开始时间
	EndDate                 time.Time         //投放结束时间
	DayParting              [][]bool          //投放时间段
	Countries               StringSet         //要投放的国家
	States                  StringSet         //要投放的地区
	Cities                  StringSet         //要投放的城市
	ADExchanges             uint32            //要投放的adx1,2,4和
	ADXisALL                bool              //是否全选adx true是
	SourceType              uint8             //投放app还是网站 网站：1 app：2 全部：0
	ControlType             uint8             //流量控制类型 BLACKLIST:1 WHITELIST:2 NONE:0
	ControlList             StringSet         //流量控制名单
	Categories              StringSet         //投放的广告类别
	ConnectionType          StringSet         // 连接类型：Unknown:0 Ethernet:1 Wifi:2 Cellular	data	– Unknown	Generation:3 Cellular	data	– 2G:4 Cellular	data	– 3G:5 Cellular	data	– 4G:6
	Carriers                StringSet         //运营商
	DeviceTypes             StringSet         //设备类型TABLET,MOBILE,1:手机/平板电脑  2:个人电脑 3:	联网电视 4:	手机 5:	平板电脑 6:	联网设备 7:	数字电视机顶盒
	OS                      StringSet         //操作系统
	IsIdfaGaid              bool              //安卓或苹果的主动接受广告选项是否打开 是：true  否：false
	IPAddress               StringSet         //投放广告的ip地址或地址段
	Audiences               StringSet         //过滤设备，名单列表
	AudienceType            uint8             //过滤设备类型 BLACKLIST:1 WHITELIST:2 NONE:0
	Retargeting             uint16            //Viewers:1,Visitors:2,Converters:4
	RedirectURL             string            //广告地址 flow，popup，Destination native
	Creatives               []Creative        //Creatives连接
	Score                   int               //得分,权重
	FilterSet               uint32            //过滤器预设置
	Adomain                 string            //广告主的主要域名或是顶级域名，用于广告主检测。对于动态物料，该参数的值可以是字符串数组。然而交易平台可以只允许一个广告主域名。
	ImpressionURL           string            //展示后的通知链接地址
	RedirectType            string            //跳转类型DESTINATION,CREATIVE,FLOW
	IsAdult                 uint8             //是否包含成人流量 1选中，0为不选 不包含成人流量
}

//根据id 获取compaign
func (cmg Compaign) GetComByID(cid uint32) Compaign {
	var compaign Compaign
	//TODO 待验证
	cmd := []byte(`{"cmd":"` + config.DSPSQL["GetComByID"] + fmt.Sprint(cid) + ` ` + ` and (adExchanges & ` + strconv.Itoa(int(config.ADXCode[config.ADX])) + `)!=0"}`)
	data, err := util.DoBytesPost(config.CompaignURL, cmd)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取compaign错误",
			"err", err,
		)
	}
	cm := cmg.GetData(data)
	if len(cm) == 1 {
		compaign = cm[0]
	}
	return compaign
}

//获取compaign
func (cmgp Compaign) GetData(data []byte) []Compaign {
	var cmps ResponseCompaign
	arrcmp := make([]Compaign, 0)
	err := json.Unmarshal(data, &cmps)
	if err != nil {
		cid, _, _, _ := jsonparser.Get(data, "data", "[0]", "campaign", "id")
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"msg", "映射compaign Struct失败。id:"+string(cid),
			// "data", string(data),
			"err", err,
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
		)
		return nil
	}
	if cmps.Status != 1 || len(cmps.Data) == 0 { //错误或没有数据
		return nil
	}
	for _, v := range cmps.Data {
		// fmt.Println(v.Campaign.ID, v.Campaign.Status)
		cmg := new(Compaign)
		//compaignid
		cmg.ID = uint32(v.Campaign.TrackingCampaignID)
		//用户id
		cmg.UID = uint32(v.Campaign.UserID)
		//时区
		cmg.TimeZone = *&TZ{
			v.Campaign.ActivityPeriodsTimezone,
			v.Campaign.ActivityPeriodsTimeArea,
		}
		//广告类型NATIVE：1,BANNER：2,POPUP：3
		switch v.Campaign.Type {
		case "NATIVE":
			cmg.Type = 1
		case "BANNER":
			cmg.Type = 2
		case "POPUP":
			cmg.Type = 3
		}
		//最大千次展示金额。千次内不超过max bid
		cmg.Maxbid = uint32(v.Campaign.BidPrice)
		// cmg.MaxbidRecords = *&MaxbidRecord{
		// 	0,
		// 	0,
		// }
		//是否不限制预算 true:不限制 false限制
		cmg.UnlimitBudget = v.Campaign.UnlimitedBudget
		//日预算
		cmg.DailyBudget = uint32(v.Campaign.DailyBudget)

		sinceTime, tillTime, _ := util.BgeinAndEndDAYOfZone(v.Campaign.ActivityPeriodsTimeArea)

		cmg.DailyBudgetRecores = *&DailyBudgetRecore{
			sinceTime,
			tillTime,
			0,
		}
		//==============设置条件
		if !cmg.UnlimitBudget {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["MaxBidFilter"]
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["DailyBudgetFilter"]
		}
		//1为ASAP，2为smooth
		switch v.Campaign.SpendStrategy {
		case "ASAP":
			cmg.SpendStrategy = 1
		case "Smooth":
			cmg.SpendStrategy = 2
		}
		//==================设置条件
		if cmg.SpendStrategy == 2 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["SpendStrategyFilter"]
		}
		//每天每个广告位的预算
		cmg.DailyPerPlacementBudget = uint32(v.Campaign.DailyPerPlacementBudget)
		cmg.DailyPPBRecords = *&DailyPPBRecord{
			make(map[string]uint32, config.MapLength),
		}
		//==================设置条件
		if cmg.DailyPerPlacementBudget > 0 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["DailyPPBFilter"]
		}
		//总预算
		cmg.TotalBudget = uint32(v.Campaign.TotalBudget)
		cmg.TotalBudgetRecords = *&TotalBudgetRecord{
			0,
		}
		//==================设置条件
		if cmg.TotalBudget > 0 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["TotalBudgetFilter"]
		}
		//是否开启投放频率
		cmg.FreqCapEnabled = v.Campaign.FreqCapEnabled
		//==================设置条件
		if cmg.FreqCapEnabled {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["FreqCapFilter"]
		}
		//DEVICE-设备:1，USER-用户:2
		switch v.Campaign.FreqCapType {
		case "DEVICE":
			cmg.FreqCapType = 1
		case "USER":
			cmg.FreqCapType = 2
		}

		//投放次数限制
		if cmg.FreqCapEnabled { //如果开启投放频率控制
			cmg.FreqCountLimit = uint32(v.Campaign.FreqCountLimit)
			//3,6,12,24小时
			cmg.FreqTimeWindow = uint8(v.Campaign.FreqTimeWindow)
			//投放记录,如用户更改投放策略则数据清空
			sinceTime = time.Now().UTC()
			interval := fmt.Sprint(cmg.FreqTimeWindow) + "h"
			d, _ := time.ParseDuration(interval)
			tillTime = sinceTime.Add(d)
			// pp.Println("compaign", int(cmg.ID), sinceTime, tillTime)
			cmg.FreqRecords = *&FreqRecord{
				sinceTime,
				tillTime,
				make(map[string]uint32, config.MapLength),
				make(map[string]uint32, config.MapLength),
			}
		}

		//投放开始时间
		loc, _ := time.LoadLocation(v.Campaign.ActivityPeriodsTimeArea)
		cmg.StartDate, _ = time.ParseInLocation("2006-01-02 15:04", v.Campaign.FromDate+" "+v.Campaign.FromTime, loc)
		// cmg.StartDate, _ = time.Parse("2006-01-02 15:04 -0700", v.Campaign.FromDate+" "+v.Campaign.FromTime+" "+strings.Replace(v.Campaign.ActivityPeriodsTimezone, ":", "", -1))
		//投放结束时间
		cmg.EndDate, _ = time.ParseInLocation("2006-01-02 15:04", v.Campaign.ToDate+" "+v.Campaign.ToTime, loc)
		// cmg.EndDate, _ = time.Parse("2006-01-02 15:04 -0700", v.Campaign.ToDate+" "+v.Campaign.ToTime+" "+strings.Replace(v.Campaign.ActivityPeriodsTimezone, ":", "", -1))
		//==================设置条件
		if v.Campaign.FromDate != "" && v.Campaign.ToDate != "" {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["BidderTimeFilter"]
		}
		//投放时间段
		cmg.DayParting = v.Campaign.DayParting
		//要投放的国家
		bm1 := roaring.NewBitmap()
		for _, vs := range v.Campaign.Countries {
			bm1.Add(util.Hashcode(vs))
		}
		var isall bool
		if v.Campaign.CountriesIsAll == 0 {
			isall = false
		} else {
			isall = true
		}
		cmg.Countries = *&StringSet{
			v.Campaign.Countries,
			bm1,
			isall,
		}
		//==================设置条件
		if !cmg.Countries.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["CountriesFilter"]
		}
		//要投放的地区
		bm2 := roaring.NewBitmap()
		for _, vs := range v.Campaign.States {
			bm2.Add(util.Hashcode(vs))
		}
		if len(v.Campaign.States) == 0 {
			isall = true
		} else {
			isall = false
		}

		cmg.States = *&StringSet{
			v.Campaign.States,
			bm2,
			isall,
		}
		//==================设置条件
		if !cmg.States.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["StatesFilter"]
		}
		//要投放的city
		bm3 := roaring.NewBitmap()
		for _, vs := range v.Campaign.Cities {
			bm3.Add(util.Hashcode(vs))
		}
		if len(v.Campaign.Cities) == 0 {
			isall = true
		} else {
			isall = false
		}
		cmg.Cities = *&StringSet{
			v.Campaign.Countries,
			bm3,
			isall,
		}
		//==================设置条件
		if !cmg.Cities.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["CitiesFilter"]
		}
		//要投放的adx
		switch v.Campaign.AdxIsAll {
		case 1:
			cmg.ADXisALL = true
		case 0:
			cmg.ADXisALL = false
		}

		cmg.ADExchanges = uint32(v.Campaign.AdExchanges)
		//==================设置条件
		// if !cmg.ADXisALL {
		// 	cmg.FilterSet = cmg.FilterSet + config.FilterCode["ADXFilter"]
		// }
		//投放app还是网站 网站：1 app：2 全部：0
		switch v.Campaign.SourceType {
		case "SITE":
			cmg.SourceType = 1
		case "APP":
			cmg.SourceType = 2
		case "ANY":
			cmg.SourceType = 0
		}
		//流量控制类型 BLACKLIST:1 WHITELIST:2 NONE 0
		switch v.Campaign.ClientType {
		case "BLACKLIST":
			cmg.ControlType = 1
		case "WHITELIST":
			cmg.ControlType = 2
		case "NONE":
			cmg.ControlType = 0
		}
		//==================设置条件
		if cmg.ControlType != 0 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["ControlListFilter"]
		}
		//	流量控制名单
		bm4 := roaring.NewBitmap()
		for _, vs := range v.Campaign.ClientIds {
			bm4.Add(util.Hashcode(vs))
		}
		cmg.ControlList = *&StringSet{
			v.Campaign.ClientIds,
			bm4,
			false,
		}
		//投放的广告类别
		bm5 := roaring.NewBitmap()
		tmpArr := make([]string, len(v.Campaign.Categories))
		for k, vs := range v.Campaign.Categories {
			bm5.Add(util.Hashcode(vs.Key))
			tmpArr[k] = vs.Key
		}
		if v.Campaign.CountriesIsAll == 0 {
			isall = false
		} else {
			isall = true
		}
		cmg.Categories = *&StringSet{
			tmpArr,
			bm5,
			isall,
		}
		//==================设置条件
		if !cmg.Categories.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["CategoriesFilter"]
		}
		// 连接类型：
		bm6 := roaring.NewBitmap()
		tmpArr = make([]string, len(v.Campaign.ConnectionType))
		for k, vs := range v.Campaign.ConnectionType {
			bm6.Add(vs)
			tmpArr[k] = fmt.Sprint(vs)
		}

		cmg.ConnectionType = *&StringSet{
			tmpArr,
			bm6,
			false,
		}
		//==================设置条件
		if len(v.Campaign.ConnectionType) < 7 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["ConnectTypeFilter"]
		}
		//运营商
		bm7 := roaring.NewBitmap()
		for _, vs := range v.Campaign.Carriers {
			bm7.Add(util.Hashcode(vs))
		}
		if len(v.Campaign.Carriers) == 0 {
			isall = true
		} else {
			isall = false
		}
		cmg.Carriers = *&StringSet{
			v.Campaign.Carriers,
			bm7,
			isall,
		}
		//==================设置条件
		if !cmg.Carriers.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["CarriersFilter"]
		}
		//设备类型TABLET,MOBILE,1:手机/平板电脑  2:个人电脑 3:	联网电视 4:	手机 5:	平板电脑 6:	联网设备 7:	数字电视机顶盒
		dtarr := make([]string, 0)
		bm8 := roaring.NewBitmap()
		for _, vs := range v.Campaign.DeviceTypes {
			bm8.Add(vs)
			dtarr = append(dtarr, fmt.Sprint(vs))
		}

		cmg.DeviceTypes = *&StringSet{
			dtarr,
			bm8,
			false,
		}
		//==================设置条件
		if len(v.Campaign.DeviceTypes) < 6 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["DeviceTypeFilter"]
		}
		//操作系统
		bm9 := roaring.NewBitmap()
		osarr := make([]string, 0)
		for _, vs := range v.Campaign.Oses {
			bm9.Add(util.Hashcode(strings.ToLower(vs)))
			osarr = append(osarr, strings.ToLower(vs))
		}
		if v.Campaign.OsIsAll == 0 {
			isall = false
		} else {
			isall = true
		}
		cmg.OS = *&StringSet{
			osarr,
			bm9,
			isall,
		}
		//==================设置条件
		if !cmg.OS.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["OSFilter"]
		}
		//安卓或苹果的主动接受广告选项是否打开 是：true-检查是否有  否：false-不必检查 设备id号为ifa
		cmg.IsIdfaGaid = v.Campaign.IsIdfaGaid
		//投放广告的ip地址或地址段
		bm10 := roaring.NewBitmap()
		tmpArr = make([]string, 0)
		for _, vs := range v.Campaign.Ips {
			if strings.Contains(vs, "/") { //是地址段
				tmpArr = append(tmpArr, vs)
			} else {
				bm10.Add(util.Hashcode(vs))
			}
		}
		if len(v.Campaign.Ips) == 0 {
			isall = true
		} else {
			isall = false
		}
		cmg.IPAddress = *&StringSet{
			tmpArr,
			bm10,
			isall,
		}
		//==================设置条件
		if !cmg.IPAddress.All {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["IPFilter"]
		}
		//过滤设备，名单列表
		bm11 := roaring.NewBitmap()
		for _, vs := range v.Campaign.AudienceIds {
			bm11.Add(util.Hashcode(vs))
		}
		cmg.Audiences = *&StringSet{
			v.Campaign.AudienceIds,
			bm11,
			false,
		}

		//过滤设备类型 BLACKLIST:1 WHITELIST:2 NONE 0
		switch v.Campaign.AudienceType {
		case "BLACKLIST":
			cmg.AudienceType = 1
		case "WHITELIST":
			cmg.AudienceType = 2
		case "NONE":
			cmg.AudienceType = 0
		}
		//==================设置条件
		if cmg.AudienceType != 0 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["AudiencesFilter"]
		}
		//Viewers:1,Visitors:2,Converters:4
		cmg.Retargeting = uint16(v.Campaign.Retarget)
		//广告地址 flow，popup，Destination
		cmg.RedirectURL = v.Campaign.RedirectURL
		//Score
		cmg.Score = int(cmg.Maxbid)
		//Creatives连接
		if cmg.Type == 1 { //nativebanner
			cmg.Creatives = make([]Creative, 0)
			for _, vs := range v.Creatives {
				if len(vs.Banner.MainImage) > 0 && vs.Status {
					creativeid := vs.UUID
					offerurl := vs.DestinationURL
					imgurl := config.CDNURL + vs.Banner.MainImage[0].BannerMeta.CdnURL
					imgheight := vs.Banner.MainImage[0].BannerMeta.Height
					imgwidth := vs.Banner.MainImage[0].BannerMeta.Width
					imgsize := vs.Banner.MainImage[0].BannerMeta.Size
					imgmime := vs.Banner.MainImage[0].BannerMeta.Mime
					iconURL := config.CDNURL + vs.Banner.IconImage[0].BannerMeta.CdnURL
					iconHeight := vs.Banner.IconImage[0].BannerMeta.Height
					iconWidth := vs.Banner.IconImage[0].BannerMeta.Width
					iconSize := vs.Banner.IconImage[0].BannerMeta.Size
					iconMime := vs.Banner.IconImage[0].BannerMeta.Mime
					headline := vs.Banner.Headline
					description := vs.Banner.Description
					buttontext := vs.Banner.ButtonText
					imgoffer := &Creative{
						creativeid,
						offerurl,
						imgurl,
						imgheight,
						imgwidth,
						imgsize,
						imgmime,
						iconURL,
						iconHeight,
						iconWidth,
						iconSize,
						iconMime,
						headline,
						description,
						buttontext,
						strconv.Itoa(imgwidth) + "x" + strconv.Itoa(imgheight),
						strconv.Itoa(iconWidth) + "x" + strconv.Itoa(iconHeight),
					}
					cmg.Creatives = append(cmg.Creatives, *imgoffer)
				}
			}
		}
		//
		if cmg.Type == 2 { //banner
			cmg.Creatives = make([]Creative, 0)
			for _, vs := range v.Creatives {
				if vs.Banner.BannerMeta.CdnURL != "" && vs.Status {
					creativeid := vs.UUID
					offerurl := vs.DestinationURL
					imgurl := config.CDNURL + vs.Banner.BannerMeta.CdnURL
					imgheight := vs.Banner.BannerMeta.Height
					imgwidth := vs.Banner.BannerMeta.Width
					imgsize := vs.Banner.BannerMeta.Size
					imgmime := vs.Banner.BannerMeta.Mime
					imgoffer := &Creative{
						creativeid,
						// util.DealXMLURL(offerurl),
						// util.DealXMLURL(imgurl),
						offerurl,
						imgurl,
						imgheight,
						imgwidth,
						imgsize,
						imgmime,
						"",
						0,
						0,
						0,
						"",
						"",
						"",
						"",
						strconv.Itoa(imgwidth) + "x" + strconv.Itoa(imgheight),
						"",
					}

					cmg.Creatives = append(cmg.Creatives, *imgoffer)
				}
			}
		}
		//==================设置条件 图片

		if cmg.Type == 1 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["NativeImageTypeFilter"]
		}
		if cmg.Type == 2 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["ImageTypeFilter"]
		}

		//==========广告主的主要域名或是顶级域名，用于广告主检测。对于动态物料，该参数的值可以是字符串数组。然而交易平台可以只允许一个广告主域名。
		cmg.Adomain = v.Campaign.Domain
		//展示后的通知链接地址
		cmg.ImpressionURL = v.Campaign.ImpressionURL
		//状态
		cmg.Status = v.Campaign.Status
		//==================设置条件
		if cmg.Status == config.CompaignStatus[config.CS_PENDING] {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["StatusFilter"]
		}
		//====================是否包含成人流量
		cmg.IsAdult = v.Campaign.ISAdult
		//==================设置条件
		if cmg.IsAdult == 0 {
			cmg.FilterSet = cmg.FilterSet + config.FilterCode["AdultFilter"]
		}
		//================================
		cmg.RedirectType = v.Campaign.RedirectType
		arrcmp = append(arrcmp, *cmg)
	}

	return arrcmp
}

//集合结构
type StringSet struct {
	Data []string        //数据
	Bm   *roaring.Bitmap //hashcode后的bitmap
	All  bool            //是否全选
}

type TZ struct {
	Offset   string
	ZoneName string
}

//总预算记录
type TotalBudgetRecord struct {
	Cost uint32 //总消费
}

//每天每个广告位的预算记录
type DailyPPBRecord struct {
	Offer map[string]uint32 //每个广告位投放数量
}

//日预算记录
type DailyBudgetRecore struct {
	SinceTime time.Time //开始计算时间
	TillTime  time.Time //结束计算时间
	Cost      uint32    //消费
}

//千次展示记录
// type MaxbidRecord struct {
// 	Count uint32 //展示次数
// 	Cost  uint32 //花费
// }

//投放记录,如用户更改投放策略则数据清空
type FreqRecord struct {
	SinceTime time.Time         //开始计算时间
	TillTime  time.Time         //结束计算时间
	User      map[string]uint32 //用户标识
	Device    map[string]uint32 //设备标识
}

//图片连接
type Creative struct {
	CreativeID    string
	Url           string //广告地址
	ImgUrl        string //主图地址
	Height        int    //高 单位pix
	Width         int    //宽 单位pix
	Size          int    //尺寸 单位Byte TODO 确认单位
	Mime          string //类型
	IconURL       string
	IconHeight    int
	IconWidth     int
	IconSize      int
	IconMime      string
	Headline      string
	Description   string
	Button        string //按钮名字，点击后跳转
	MainImgFormat string //大图尺寸格式230x34
	IconFormat    string //icon 尺寸格式58x50
}

type ResponseCompaign struct {
	Data []struct {
		Campaign struct {
			ISAdult                 uint8    `json:"adult"`
			ActivityPeriodsTimeArea string   `json:"activityPeriodsTimeArea"`
			ActivityPeriodsTimezone string   `json:"activityPeriodsTimezone"`
			AdExchanges             int      `json:"adExchanges"`
			AdxIsAll                int      `json:"adx_is_all"`
			AudienceIds             []string `json:"audienceIds"`
			AudienceType            string   `json:"audienceType"`
			BidPrice                int      `json:"bidPrice"`
			Carriers                []string `json:"carriers"`
			Categories              []struct {
				ID   int    `json:"id"`
				Key  string `json:"key"`
				Name string `json:"name"`
			} `json:"categories"`
			CategoriesIsAll         int      `json:"categories_is_all"`
			Cities                  []string `json:"cities"`
			ClientIds               []string `json:"clientIds"`
			ClientType              string   `json:"clientType"`
			ConnectionType          []uint32 `json:"connectionType"`
			ConversionActionType    string   `json:"conversionActionType"`
			ConversionActionURL     string   `json:"conversionActionUrl"`
			Countries               []string `json:"countries"`
			CountriesIsAll          int      `json:"countries_is_all"`
			DailyBudget             int      `json:"dailyBudget"`
			DailyPerPlacementBudget int      `json:"dailyPerPlacementBudget"`
			DayParting              [][]bool `json:"dayParting"`
			DeviceTypes             []uint32 `json:"deviceTypes"`
			Domain                  string   `json:"domain"`
			FlowID                  string   `json:"flowId"`
			FreqCapEnabled          bool     `json:"freqCapEnabled"`
			FreqCapType             string   `json:"freqCapType"`
			FreqCountLimit          int      `json:"freqCountLimit"`
			FreqTimeWindow          int      `json:"freqTimeWindow"`
			FromDate                string   `json:"fromDate"`
			FromTime                string   `json:"fromTime"`
			Hash                    string   `json:"hash"`
			ID                      int      `json:"id"`
			ImpressionURL           string   `json:"impressionURL"`
			Ips                     []string `json:"ips"`
			IsIdfaGaid              bool     `json:"isIdfaGaid"`
			Name                    string   `json:"name"`
			OsIsAll                 int      `json:"os_is_all"`
			Oses                    []string `json:"oses"`
			RedirectType            string   `json:"redirectType"`
			RedirectURL             string   `json:"redirectUrl"`
			Retarget                int      `json:"retarget"`
			RevenueType             string   `json:"revenueType"`
			RevenueValue            int      `json:"revenueValue"`
			Score                   int      `json:"score"`
			SourceType              string   `json:"sourceType"`
			SpendStrategy           string   `json:"spendStrategy"`
			States                  []string `json:"states"`
			Status                  uint8    `json:"status"`
			ToDate                  string   `json:"toDate"`
			ToTime                  string   `json:"toTime"`
			TotalBudget             int      `json:"totalBudget"`
			TrackingCampaignID      int      `json:"trackingCampaignId"`
			Type                    string   `json:"type"`
			UnlimitedBudget         bool     `json:"unlimitedBudget"`
			UserID                  int      `json:"userId"`
		} `json:"campaign"`
		Creatives []struct {
			ApprovalStatus string `json:"approvalStatus"`
			Banner         struct {
				Brand       string `json:"brand"`
				Description string `json:"brandingText"`
				ButtonText  string `json:"ctaText"`
				Headline    string `json:"headline"`
				BannerMeta  struct {
					CdnURL string `json:"cdnUrl"`
					EtagB  string `json:"etag"`
					Height int    `json:"height"`
					Mime   string `json:"mime"`
					Size   int    `json:"size"`
					Width  int    `json:"width"`
				} `json:"bannerMeta"`
				IconImage []struct {
					BannerMeta struct {
						CdnURL string `json:"cdnUrl"`
						EtagI  string `json:"etag"`
						Height int    `json:"height"`
						Mime   string `json:"mime"`
						Size   int    `json:"size"`
						Width  int    `json:"width"`
					} `json:"bannerMeta"`
					ID           string `json:"id"`
					OriginalName string `json:"originalName"`
				} `json:"iconImage"`
				ID        string `json:"id"`
				MainImage []struct {
					BannerMeta struct {
						CdnURL string `json:"cdnUrl"`
						EtagM  string `json:"etag"`
						Height int    `json:"height"`
						Mime   string `json:"mime"`
						Size   int    `json:"size"`
						Width  int    `json:"width"`
					} `json:"bannerMeta"`
					ID           string `json:"id"`
					OriginalName string `json:"originalName"`
				} `json:"mainImage"`
			} `json:"banner"`
			Delete           bool     `json:"delete"`
			DestinationURL   string   `json:"destinationUrl"`
			NativeBanner     struct{} `json:"nativeBanner"`
			RejectionComment string   `json:"rejectionComment"`
			Status           bool     `json:"status"`
			UUID             string   `json:"uuid"`
		} `json:"creatives"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
