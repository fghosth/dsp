package model

import (
	"jvole.com/dsp/config"
)

type Adsterra struct {
	ADX ADXbase
}

func NewAdsterraADX(data []byte, code uint32, domain string) *Adsterra {
	adx := &Adsterra{
		*NewADXbase(data, code, domain),
	}
	return adx
}

//获得竞价请求的字符串
func (st *Adsterra) GetBidderResponse(cp Compaign) (response string) {
	switch st.GetType() { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		// response = st.NativeResponse(cp, config.BASE_BID_TEMPLATE, config.MGID_NATIVE_JSON_TEMPLATE)
	case 2:
		// response = st.BannerResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_BANNER_XML_TEMPLATE)
	case 3:
		response = st.PopupResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_POPUP_XML_TEMPLATE)
	}
	return
}

//GetUA()
func (st *Adsterra) GetUA() string {
	return st.ADX.GetUA()
}

//获得site的page
func (st *Adsterra) GetSitePage() string {
	return st.ADX.GetSitePage()
}

//获得site的IDzone
func (st *Adsterra) GetSiteIDzone() int {
	return st.ADX.GetSiteIDzone()
}

//获得language
func (st *Adsterra) GetLanguage() string {
	return st.ADX.GetLanguage()
}

//获得osversion
func (st *Adsterra) GetOSVStr() string {
	return st.ADX.GetOSVStr()
}

//获得os
func (st *Adsterra) GetOSStr() string {
	return st.ADX.GetOSStr()
}

//获得model
func (st *Adsterra) GetModel() string {
	return st.ADX.GetModel()
}

//获得广告位图片尺寸
func (st *Adsterra) GetImages() []string {
	return st.ADX.GetImages()
}

func (st *Adsterra) GetOfferInfo() OfferInfo {
	return st.ADX.GetOfferInfo()
}

func (st *Adsterra) GetData() []byte {
	return st.ADX.GetData()
}

func (st *Adsterra) SetData(data []byte) {
	st.ADX.data = data
}

//设置offerinfo
func (st *Adsterra) SetOfferInfo() {
	st.ADX.SetOfferInfo()
	st.ADX.offerInfo.Images = st.GetImages()
}

//获得广告唯一id
func (st *Adsterra) GetOfferID() string {
	return st.ADX.GetOfferID()
}

//获得impid
func (st *Adsterra) GetIMPID() string {
	return st.ADX.GetIMPID()
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *Adsterra) GetType() uint32 {
	return st.ADX.GetType()
}

//获得低价 *1000000
func (st *Adsterra) GetBidFloor() uint32 {
	return st.ADX.GetBidFloor()
}

//获得广告位标识
func (st *Adsterra) GetPostion() string {

	return st.ADX.GetPostion()
}

//获得设备号
func (st *Adsterra) GetDeviceID() string {
	return st.ADX.GetDeviceID()
}

//获得用户标识
func (st *Adsterra) GetUserID() string {
	return st.ADX.GetUserID()
}

//获得country
func (st *Adsterra) GetCountry() string {
	return st.ADX.GetCountry()
}

//获得city
func (st *Adsterra) GetCity() string {
	return st.ADX.GetCity()
}

//获得region/state
func (st *Adsterra) GetRegion() string {
	return st.ADX.GetRegion()
}

//获得ADX
func (st *Adsterra) GetADX() uint32 {
	return st.ADX.GetADX()
}

//获得ADX
func (st *Adsterra) SetCode(code uint32) {
	st.ADX.SetCode(code)
}

//获得App或site标识
func (st *Adsterra) GetAppSite() string {
	return st.ADX.GetAppSite()
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *Adsterra) GetSourceType() uint32 {

	return st.ADX.GetSourceType()
}

//获得广告类型
func (st *Adsterra) GetContentCat() []string {

	return st.ADX.GetContentCat()
}

//获得连接类型
func (st *Adsterra) GetConnType() uint32 {
	return st.ADX.GetConnType()
}

//获得运营商
func (st *Adsterra) GetCarriers() string {
	return st.ADX.GetCarriers()
}

//获得设备类型
func (st *Adsterra) GetDeviceType() uint32 {

	return st.ADX.GetDeviceType()
}

//获得os系统
func (st *Adsterra) GetOS() string {
	return st.ADX.GetOS()
}

//获得IdfaGaid
func (st *Adsterra) GetIdfaGaid() string {
	return st.ADX.GetIdfaGaid()
}

//获得IP
func (st *Adsterra) GetIP() string {
	return st.ADX.GetIP()
}

//获得native的response字符串 NativeResponse
func (st *Adsterra) NativeResponse(cp Compaign, template, admTemplate string) string {
	return ""
}

//获得banner的response字符串
func (st *Adsterra) BannerResponse(cp Compaign, template, admTemplate string) string {

	return ""
}

//获得Popup的response字符串
func (st *Adsterra) PopupResponse(cp Compaign, template, admTemplate string) string {

	return st.ADX.PopupResponse(cp, template, admTemplate)
}
