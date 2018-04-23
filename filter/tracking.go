package filter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

func trackingLevelFilter(cid uint32, v interface{}, offer model.Offer) {
	objtype := reflect.TypeOf(v).String()
	var cmpvalue, ofvalue interface{}
	compaign := index.CPINDEX.GetCompaign(cid)
	// tf := &TypeFilter{1}
	// sf := &SourceFilter{2}
	//
	// igf := &IdfaGaidFilter{16}
	// levelf := []LevelFilter{igf, tf, sf}
	switch objtype {

	case "*filter.TypeFilter":
		cmpvalue = fmt.Sprint("广告类型NATIVE：1,BANNER：2,POPUP：3", compaign.Type)
		ofvalue = offer.GetType()
	case "*filter.SourceFilter":
		cmpvalue = fmt.Sprint("投放app还是网站 网站：1 app：2 全部：0", compaign.SourceType)
		ofvalue = offer.GetSourceType()
	case "*filter.IdfaGaidFilter":
		cmpvalue = compaign.IsIdfaGaid
		ofvalue = offer.GetIdfaGaid()
	}
	level.Debug(*logger).Log(
		"compaignID", cid,
		"msg", objtype+"==========Level过滤掉此compaign===========",
		"compaign", cmpvalue,
		"adx/records:", ofvalue,
	)
	str_content := fmt.Sprint("compaignID", cid, "msg", objtype+"Level过滤掉此compaign:", "compaign", cmpvalue, "adx/records:", ofvalue)
	util.Tracefile(str_content, config.TrackingFile)
}

func trackingCPFilter(v interface{}, offer model.Offer, compaign model.Compaign) {
	cid := compaign.ID
	objtype := reflect.TypeOf(v).String()
	var cmpvalue, ofvalue interface{}
	switch objtype {
	case "*filter.AdultFilter":
		cmpvalue = fmt.Sprint("是否接受adult流量:接受：1 不接受：0 compaign：", compaign.IsAdult)
		ofvalue = offer.GetOfferInfo().Domain_STR
	case "*filter.MaxBidFilter":
		cmpvalue = compaign.Maxbid
		ofvalue = offer.GetBidFloor()
	case "*filter.DailyBudgetFilter":
		cmpvalue = compaign.DailyBudget
		ofvalue = compaign.DailyBudgetRecores.Cost
	case "*filter.SpendStrategyFilter":
		cmpvalue = compaign.SpendStrategy
		ofvalue = compaign.DailyBudgetRecores
	case "*filter.TotalBudgetFilter":
		cmpvalue = compaign.TotalBudget
		ofvalue = compaign.TotalBudgetRecords.Cost
	case "*filter.FreqCapFilter":
		cmpvalue = fmt.Sprint("类型:1设备，2用户:", compaign.FreqCapType, "次数限制:", compaign.FreqCountLimit, "时间段:", compaign.FreqTimeWindow)
		ofvalue = compaign.FreqRecords
	case "*filter.BidderTimeFilter":
		res := util.ZoneOffset(compaign.TimeZone.Offset)
		zoneName := compaign.TimeZone.ZoneName
		//TODO 有可能错误
		rest := time.Now().UTC().In(time.FixedZone(zoneName, res))
		cmpvalue = fmt.Sprint(compaign.StartDate, compaign.EndDate)
		ofvalue = rest
	case "*filter.DayPartingFilter":
		res := util.ZoneOffset(compaign.TimeZone.Offset)
		// fmt.Println("time", res)
		zoneName := compaign.TimeZone.ZoneName
		rest := time.Now().UTC().In(time.FixedZone(zoneName, res))
		week := int(rest.Weekday()) //获得周0，1，2，3，4，5，6 sun,mon。。。。。
		hour := rest.Hour()
		cmpvalue = compaign.DayParting
		ofvalue = fmt.Sprint("week:", week, "hour:", hour)
	case "*filter.CountriesFilter":
		cmpvalue = compaign.Countries.Data
		ofvalue = offer.GetCountry()
	case "*filter.StatesFilter":
		cmpvalue = compaign.States.Data
		ofvalue = offer.GetRegion()
	case "*filter.CitiesFilter":
		cmpvalue = compaign.Cities.Data
		ofvalue = offer.GetCity()
	case "*filter.ADXFilter":
		cmpvalue = compaign.ADExchanges
		ofvalue = config.ADX
	case "*filter.ControlListFilter":
		cmpvalue = compaign.ControlList.Data
		ofvalue = offer.GetAppSite()
	case "*filter.CategoriesFilter":
		cmpvalue = compaign.Categories.Data
		ofvalue = offer.GetContentCat()
	case "*filter.CarriersFilter":
		cmpvalue = compaign.Carriers.Data
		ofvalue = offer.GetCarriers()
	case "*filter.OSFilter":
		cmpvalue = compaign.OS.Data
		ofvalue = offer.GetOS()
	case "*filter.AudiencesFilter":
		cmpvalue = compaign.Audiences.Data
		ofvalue = offer.GetIdfaGaid()
	case "*filter.DailyPPBFilter":
		cmpvalue = compaign.DailyPPBRecords.Offer[offer.GetPostion()]
		ofvalue = compaign.DailyPerPlacementBudget
	case "*filter.IPFilter":
		cmpvalue = compaign.IPAddress.Data
		ofvalue = offer.GetIP()
	case "*filter.ConnectTypeFilter":
		cmpvalue = compaign.ConnectionType.Data
		ofvalue = offer.GetConnType()
	case "*filter.DeviceTypeFilter":
		cmpvalue = compaign.DeviceTypes.Data
		ofvalue = offer.GetDeviceType()
	case "*filter.ImageTypeFilter":
		cmpvalue = compaign.Creatives
		ofvalue = offer.GetImages()
	case "*filter.NativeImageTypeFilter":
		cmpvalue = compaign.Creatives
		ofvalue = offer.GetImages()
	}
	level.Debug(*logger).Log(
		"compaignID", cid,
		"msg", objtype+"==========CP过滤掉此compaign===========",
		"compaign", cmpvalue,
		"adx/records:", ofvalue,
	)
	str_content := fmt.Sprint("compaignID", cid, "msg", objtype+"Level过滤掉此compaign", "compaign", cmpvalue, "adx/records:", ofvalue)
	util.Tracefile(str_content, config.TrackingFile)
}
