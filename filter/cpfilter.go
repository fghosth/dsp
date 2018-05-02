package filter

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"
	"jvole.com/dsp/config"
	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

type CPFilter interface { //过滤器接口 大范围过滤
	/*
		过滤器
		@param index.IndexMap
		@param offer model.Offer adx的offer
		@return Compaignid 的集合  过滤后剩下的compaign
	*/
	Filter(compaign model.Compaign, offer model.Offer) bool
	//获得过滤器编号2^n(n>=0)
	GetCode() uint32
}

//执行所有CPFilter 返回单个compaign是否可以通过过滤
func CPFiltercomp(compaign model.Compaign, offer model.Offer) bool {
	mbf := &MaxBidFilter{config.FilterCode["MaxBidFilter"]}
	dbf := &DailyBudgetFilter{config.FilterCode["DailyBudgetFilter"]}
	ssf := &SpendStrategyFilter{config.FilterCode["SpendStrategyFilter"]}
	tbf := &TotalBudgetFilter{config.FilterCode["TotalBudgetFilter"]}
	fcf := &FreqCapFilter{new(sync.RWMutex), config.FilterCode["FreqCapFilter"]}
	btf := &BidderTimeFilter{config.FilterCode["BidderTimeFilter"]}
	dpf := &DayPartingFilter{config.FilterCode["DayPartingFilter"]}
	cf := &CountriesFilter{config.FilterCode["CountriesFilter"]}
	sf := &StatesFilter{config.FilterCode["StatesFilter"]}
	cityf := &CitiesFilter{config.FilterCode["CitiesFilter"]}
	adxf := &ADXFilter{config.FilterCode["ADXFilter"]}
	clf := &ControlListFilter{config.FilterCode["ControlListFilter"]}
	catef := &CategoriesFilter{config.FilterCode["CategoriesFilter"]}
	caarf := &CarriersFilter{config.FilterCode["CarriersFilter"]}
	osf := &OSFilter{config.FilterCode["OSFilter"]}
	af := &AudiencesFilter{config.FilterCode["AudiencesFilter"]}
	dppbf := &DailyPPBFilter{new(sync.RWMutex), config.FilterCode["DailyPPBFilter"]}
	ipf := &IPFilter{config.FilterCode["IPFilter"]}
	ctf := &ConnectTypeFilter{config.FilterCode["ConnectTypeFilter"]}
	dtf := &DeviceTypeFilter{config.FilterCode["DeviceTypeFilter"]}
	itf := &ImageTypeFilter{config.FilterCode["ImageTypeFilter"]}
	statusf := &StatusFilter{config.FilterCode["StatusFilter"]}
	nitf := &NativeImageTypeFilter{config.FilterCode["StatusFilter"]}
	adultf := &AdultFilter{config.FilterCode["AdultFilter"]}

	cplf := []CPFilter{adultf, statusf, mbf, dbf, ssf, tbf, fcf, btf, dpf, cf, sf, cityf, adxf, clf, catef, caarf, osf, af, dppbf, ipf, ctf, dtf, itf, nitf}
	andf := &AndCPFilter{}
	andf.AndFilter(cplf...)
	return andf.Filter(compaign, offer)
}

//成人流量过滤。
type AdultFilter struct {
	Code uint32
}

func (adultf *AdultFilter) GetCode() uint32 {
	return adultf.Code
}

func (adultf *AdultFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if config.AdultBitmap.Contains(offerinfo.Domain) && compaign.IsAdult == 1 {
		return true
	}
	if !config.AdultBitmap.Contains(offerinfo.Domain) && compaign.IsAdult == 2 {
		return true
	}
	if compaign.IsAdult == 3 {
		return true
	}
	return false
}

//状态过滤。 pending不参与竞价
type StatusFilter struct { //编号1
	Code uint32
}

func (sf *StatusFilter) GetCode() uint32 {
	return sf.Code
}

func (sf *StatusFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	if compaign.Status == config.CompaignStatus[config.CS_PENDING] {
		return false
	}
	return true
}

//出价过滤。 floorbid 比较
type MaxBidFilter struct { //编号1
	Code uint32
}

func (mbf *MaxBidFilter) GetCode() uint32 {
	return mbf.Code
}

func (mbf *MaxBidFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.Maxbid/config.CPMBidFix < offerinfo.Bidfloor {
		return false
	}
	return true
}

//日预算
type DailyBudgetFilter struct { //编号2
	Code uint32
}

func (dbf *DailyBudgetFilter) GetCode() uint32 {
	return dbf.Code
}
func (dbf *DailyBudgetFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	if compaign.DailyBudget < compaign.DailyBudgetRecores.Cost {
		return false
	}
	return true
}

//1为ASAP，2为smooth 投放策略 1为快速投放 2为 日预算平摊到24小时
type SpendStrategyFilter struct { //编号4
	Code uint32
}

func (ssf *SpendStrategyFilter) GetCode() uint32 {
	return ssf.Code
}
func (ssf *SpendStrategyFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	if compaign.SpendStrategy == 1 { //快速投放
		return true
	}

	var perh float32 //每小时费用
	perh = float32(compaign.DailyBudget) / 24 / 60

	//当前第几个小时
	dt, _ := util.CovnNOWUTC2Location(compaign.TimeZone.ZoneName)
	minute := dt.Hour()*60 + dt.Minute()

	// pp.Println(dt, dt.Hour(), dt.Minute(), int(compaign.DailyBudgetRecores.Cost), perh*float32(minute))
	if float32(compaign.DailyBudgetRecores.Cost) >= perh*float32(minute) {
		return false
	}
	return true
}

//总预算过滤
type TotalBudgetFilter struct { //编号8
	Code uint32
}

func (tbf *TotalBudgetFilter) GetCode() uint32 {
	return tbf.Code
}
func (tbf *TotalBudgetFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	if compaign.TotalBudgetRecords.Cost >= compaign.TotalBudget {
		return false
	}
	return true
}

//投放频率过滤
type FreqCapFilter struct { //编号16
	lock *sync.RWMutex
	Code uint32
}

func (fcf *FreqCapFilter) GetCode() uint32 {
	return fcf.Code
}
func (fcf *FreqCapFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	fcf.lock.RLock()
	defer fcf.lock.RUnlock()
	switch compaign.FreqCapType {
	case 1:
		if compaign.FreqRecords.Device[offerinfo.DeviceID_STR] >= compaign.FreqCountLimit {
			return false
		}
	case 2:
		if compaign.FreqRecords.User[offerinfo.UserID_STR] >= compaign.FreqCountLimit {
			return false
		}
	}

	return true
}

//投放时间过滤
type BidderTimeFilter struct { //编号16
	Code uint32
}

func (btf *BidderTimeFilter) GetCode() uint32 {
	return btf.Code
}
func (btf *BidderTimeFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	rest := time.Now().UTC()
	if util.IsBetweenTime(compaign.StartDate, compaign.EndDate, rest) {
		return true
	} else {
		return false
	}

}

//投放时间段过滤
type DayPartingFilter struct { //编号32 2^5
	Code uint32
}

func (dpf *DayPartingFilter) GetCode() uint32 {
	return dpf.Code
}
func (dpf *DayPartingFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	rest, _ := util.CovnNOWUTC2Location(compaign.TimeZone.ZoneName)
	week := int(rest.Weekday()) //获得周0，1，2，3，4，5，6 sun,mon。。。。。
	hour := rest.Hour()
	if compaign.DayParting[week][hour] {
		return true
	} else {
		return false
	}
}

//要投放的国家
type CountriesFilter struct { //编号2^6
	Code uint32
}

func (cf *CountriesFilter) GetCode() uint32 {
	return cf.Code
}
func (cf *CountriesFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if offerinfo.Country == 0 { //如果没有则返回成功
		return true
	}
	if compaign.Countries.Bm.Contains(offerinfo.Country) {
		return true
	} else {
		return false
	}
}

//要投放的地区
type StatesFilter struct { //编号2^7
	Code uint32
}

func (sf *StatesFilter) GetCode() uint32 {
	return sf.Code
}
func (sf *StatesFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if offerinfo.Country == 0 { //如果没有则返回成功
		return true
	}
	if compaign.States.Bm.Contains(offerinfo.Region) {
		return true
	} else {
		return false
	}
}

//要投放的城市
type CitiesFilter struct { //编号2^8
	Code uint32
}

func (cityf *CitiesFilter) GetCode() uint32 {
	return cityf.Code
}
func (cityf *CitiesFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if offerinfo.Country == 0 { //如果没有则返回成功
		return true
	}
	if compaign.Cities.Bm.Contains(offerinfo.City) {
		return true
	} else {
		return false
	}
}

//要投放的ADX
type ADXFilter struct { //编号2^8
	Code uint32
}

func (adxf *ADXFilter) GetCode() uint32 {
	return adxf.Code
}
func (adxf *ADXFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	if (compaign.ADExchanges & offer.GetADX()) == 0 {
		return false
	} else {
		return true
	}
}

//流量黑白名单
type ControlListFilter struct { //编号2^9
	Code uint32
}

func (clf *ControlListFilter) GetCode() uint32 {
	return clf.Code
}
func (clf *ControlListFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.ControlType == 1 { //黑名单
		if compaign.ControlList.Bm.Contains(offerinfo.APPSite) {
			return false
		} else {
			return true
		}
	} else if compaign.ControlType == 2 { //白名单
		if compaign.ControlList.Bm.Contains(offerinfo.APPSite) {
			return true
		} else {
			return false
		}
	} else { //不筛选
		return true
	}
}

//投放的广告类别
type CategoriesFilter struct { //编号2^10
	Code uint32
}

func (catef *CategoriesFilter) GetCode() uint32 {
	return catef.Code
}
func (catef *CategoriesFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	cate := roaring.NewBitmap()
	cate.AddMany(offerinfo.ContentCat)
	cate.And(compaign.Categories.Bm) //如果有交集说明在范围内
	if cate.GetCardinality() > 0 {
		return true
	} else {
		return false
	}
}

//运营商
type CarriersFilter struct { //编号2^11
	Code uint32
}

func (carrf *CarriersFilter) GetCode() uint32 {
	return carrf.Code
}
func (carrf *CarriersFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.Carriers.Bm.Contains(offerinfo.Carriers) {
		return true
	} else {
		return false
	}
}

//操作系统
type OSFilter struct { //编号2^12
	Code uint32
}

func (osf *OSFilter) GetCode() uint32 {
	return osf.Code
}
func (osf *OSFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.OS.Bm.Contains(offerinfo.OS) {
		return true
	} else {
		return false
	}
}

//过滤设备，名单列表
type AudiencesFilter struct { //编号2^13
	Code uint32
}

func (af *AudiencesFilter) GetCode() uint32 {
	return af.Code
}
func (af *AudiencesFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.Audiences.Bm.Contains(offerinfo.IdfaGaid) {
		return true
	} else {
		return false
	}
}

//每天每个广告位的预算
type DailyPPBFilter struct { //编号2^14
	lock *sync.RWMutex
	Code uint32
}

func (dppbf *DailyPPBFilter) GetCode() uint32 {
	return dppbf.Code
}
func (dppbf *DailyPPBFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	dppbf.lock.RLock()
	defer dppbf.lock.RUnlock()
	_, ok := compaign.DailyPPBRecords.Offer[offerinfo.Postion]
	if !ok { //如果广告位没被记录过
		return true
	}
	if compaign.DailyPPBRecords.Offer[offerinfo.Postion] >= compaign.DailyPerPlacementBudget {
		return false
	} else {
		return true
	}
}

//IP过滤
type IPFilter struct { //编号2^15
	Code uint32
}

func (ip *IPFilter) GetCode() uint32 {
	return ip.Code
}
func (ip *IPFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()

	if compaign.IPAddress.Bm.Contains(offerinfo.IP) {
		return true
	}
	for _, v := range compaign.IPAddress.Data { //ip段
		if util.IsInsegment(offerinfo.IP_STR, v) {
			return true
		}
	}
	return false
}

//连接类型ConnectTypeFilter
type ConnectTypeFilter struct { //编号2^15
	Code uint32
}

func (ctf *ConnectTypeFilter) GetCode() uint32 {
	return ctf.Code
}
func (ctf *ConnectTypeFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if compaign.ConnectionType.Bm.Contains(offerinfo.ConnType) {
		return true
	}
	return false
}

//设备过滤 DeviceTypeFilter
type DeviceTypeFilter struct { //编号2^15
	Code uint32
}

func (dtf *DeviceTypeFilter) GetCode() uint32 {
	return dtf.Code
}
func (dtf *DeviceTypeFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	if offerinfo.DeviceType == 0 {
		return true
	}
	if compaign.DeviceTypes.Bm.Contains(offerinfo.DeviceType) {
		return true
	}
	return false
}

//banner主图片尺寸过滤 ImageTypeFilter
type ImageTypeFilter struct { //编号2^15
	Code uint32
}

func (itf *ImageTypeFilter) GetCode() uint32 {
	return itf.Code
}
func (itf *ImageTypeFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	for _, v := range compaign.Creatives {
		if v.MainImgFormat == offerinfo.Images[0] {
			return true
		}
	}
	return false
}

//native图片尺寸过滤 NativeImageTypeFilter
type NativeImageTypeFilter struct { //编号2^15
	Code uint32
}

func (nitf *NativeImageTypeFilter) GetCode() uint32 {
	return nitf.Code
}
func (nitf *NativeImageTypeFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	offerinfo := offer.GetOfferInfo()
	for _, v := range compaign.Creatives {
		imgsize := strings.Split(offerinfo.Images[0], "x")
		w, _ := strconv.Atoi(imgsize[0])
		h, _ := strconv.Atoi(imgsize[1])
		if v.Width >= w && v.Height >= h {
			return true
		}
	}
	return false
}

//=================================================
type AndCPFilter struct {
	filter []CPFilter //过滤器集合
	Code   uint32
}

func (af *AndCPFilter) GetCode() uint32 {
	return af.Code
}
func (af *AndCPFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	com := true
	for _, v := range af.filter { //遍历过滤器
		// fmt.Println("执行条件数：", k)
		if (v.GetCode() & compaign.FilterSet) == 0 { //如果不必要验证则不验证
			continue
		}
		if !v.Filter(compaign, offer) {
			if config.TrackCompaignISON { //追踪代码
				for _, cid := range config.DspTrackCIDs {
					cm := index.CPINDEX.GetCompaign(cid)
					if cid == compaign.ID && cm.ID != 0 {
						trackingCPFilter(v, offer, compaign)
					}
				}
			}
			com = false
			break
		}
	}
	return com
}

/*
多个过滤器and
*/
func (af *AndCPFilter) AndFilter(filter ...CPFilter) {
	af.filter = filter
}

type OrCPFilter struct {
	filter []CPFilter //过滤器集合
	Code   uint32
}

func (orf *OrCPFilter) GetCode() uint32 {
	return orf.Code
}

func (orf *OrCPFilter) Filter(compaign model.Compaign, offer model.Offer) bool {
	com := false
	for _, v := range orf.filter { //遍历过滤器
		if v.Filter(compaign, offer) {
			com = true
			break
		}
	}
	return com
}

/*
多个过滤器or
*/
func (af *OrCPFilter) OrFilter(filter ...CPFilter) {
	af.filter = filter
}
