package model

//adx 的结构
type Offer interface {
	//获得广告唯一id
	GetOfferID() string
	//广告类型NATIVE：1,BANNER：2,POPUP：3
	GetType() uint32
	//获得低价 *1000000
	GetBidFloor() uint32
	//获得广告位标识
	GetPostion() string
	//获得设备号
	GetDeviceID() string
	//获得用户标识
	GetUserID() string
	//获得country
	GetCountry() string
	//获得city
	GetCity() string
	//获得region/state
	GetRegion() string
	//获得ADX code
	GetADX() uint32
	//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
	GetSourceType() uint32
	//获得广告类型
	GetContentCat() []string
	//获得连接类型
	GetConnType() uint32
	//获得运营商
	GetCarriers() string
	//获得设备类型
	GetDeviceType() uint32
	//获得os系统
	GetOS() string
	//获得app或site标识
	GetAppSite() string
	//获得IdfaGaid
	GetIdfaGaid() string
	//获得IP
	GetIP() string
	//设置offer信息
	SetOfferInfo()
	//获取offer关键信息
	GetOfferInfo() OfferInfo
	//获得impid
	GetIMPID() string
	//获得竞价请求的字符串
	GetBidderResponse(compaign Compaign) string
	//获得广告位json完整数据
	GetData() []byte
	//设置广告位json完整数据
	SetData([]byte)
	//设置CODE
	SetCode(code uint32)
	//获得广告位图片尺寸
	GetImages() []string
	//GetUA()
	GetUA() string
	//获得site的page
	GetSitePage() string
	//获得site的IDzone
	GetSiteIDzone() int
	//获得language
	GetLanguage() string
	//获得osversion
	GetOSVStr() string
	//获得os
	GetOSStr() string
	//获得model
	GetModel() string
}
